package hero

import (
	"bytes"
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/config"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/payments"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	Signup(ctx context.Context, body *SignupBody) (*SignUpResult, error)
	Login(ctx context.Context, body *LoginBody) (*LoginResult, error)
	GetMagicLink(ctx context.Context, to string) error
	ValidateMagicLink(ctx context.Context, token string) (*LoginResult, error)
	ValidateJWTToken(ctx context.Context, token string) (*User, error)

	GetProfile(ctx context.Context, userId int) (*Profile, error)
	UpdateProfile(ctx context.Context, userId int, profile UpdateProfileBody) error

	GetUserAccessibleProducts(ctx context.Context, userId int) ([]ProductCard, error)
	GetUserAccessibleProduct(ctx context.Context, courseSlug string, userId int) (*ProductInfo, error)

	GetUserAccessibleLesson(ctx context.Context, lessonSlug string, userId int) (*LessonInfo, error)
	CompleteLesson(ctx context.Context, lessonSlug string, userId int) error

	ChangeAvatar(ctx context.Context, userId int, avatarPath string, avatarFileName string) error
	ChangePassword(ctx context.Context, userId int, password string) error

	GetSolvedQuizzesForQuiz(ctx context.Context, lessonSlug string, skip int, limit int) ([]QuizSolvedInfo, error)
	GetSolvedQuizzesForProduct(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error)
	GetSolvedQuizzesForUser(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error)
	SolveQuiz(ctx context.Context, dto SolveQuizDTO) error
	GetQuizBySlug(ctx context.Context, slug string) (*Quiz, error)
	DeleteSolvedQuiz(ctx context.Context, quizSlug string, userId int) error

	GetOfferForRegistration(ctx context.Context, offerSlug string) (*OfferForRegistration, error)
	GetOfferForProcessing(ctx context.Context, offerSlug string) (*OfferForProcessing, error)
	ProcessOffer(ctx context.Context, dto ProcessOfferDTO) (*ProcessOfferResult, error)

	ProcessTinkoffWebhook(ctx context.Context, payload TinkoffWebhookBody) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

const (
	MediaStatusUploaded        = "uploaded"
	RelatedMediaTypeSolvedQuiz = "solved_quiz"
)

type Service struct {
	repo   Storage
	config *config.Config
	emails IEmailsService
}

func (s *Service) GetOfferForRegistration(ctx context.Context, offerSlug string) (*OfferForRegistration, error) {
	offer, err := s.repo.GetOfferForRegistration(ctx, offerSlug)
	if err != nil && errors.Is(err, common.ErrOfferNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, common.ErrInternalError
	}

	return offer, nil
}

func (s *Service) GetOfferForProcessing(ctx context.Context, offerSlug string) (*OfferForProcessing, error) {
	offer, err := s.repo.GetOfferForProcessing(ctx, offerSlug)
	if err != nil && errors.Is(err, common.ErrOfferNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, common.ErrInternalError
	}

	return offer, nil
}

func (s *Service) ProcessOffer(ctx context.Context, dto ProcessOfferDTO) (*ProcessOfferResult, error) {
	// Шаг 1. Получить оффер со всеми полями
	offer, err := s.repo.GetOfferForProcessing(ctx, dto.Slug)
	if err != nil {
		logger.Info(ctx, err.Error())
		return nil, common.ErrInternalError
	}

	logger.Info(ctx, "got offer", "offer_id", offer.ID)

	payMethod, err := s.repo.GetPayMethod(ctx, dto.SelectedPayMethod, offer.ProjectID)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, common.ErrInternalError
	}

	logger.Info(ctx, "found pay method", "pay_method_id", payMethod.ID)

	offer.PayMethod = payMethod

	// Шаг 2. Зарегистрировать пользователя
	userId, _, err := s.createUser(ctx, CreateUserDTO{
		FirstName: dto.FirstName,
		Email:     dto.Email,
	})

	if err != nil {
		logger.Log.ErrorContext(ctx, err.Error())
		return nil, common.ErrInternalError
	}

	logger.Info(ctx, "created user", "user_id", userId)

	dto.UserID = userId

	result := ProcessOfferResult{}

	// Шаг 3. Обновить информацию о пользователе
	err = s.repo.UpdateUserInfo(ctx, UpdateUserInfoDTO{
		UserID:    userId,
		Phone:     dto.Phone,
		Telegram:  dto.Telegram,
		Instagram: dto.Instagram,
	})
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, common.ErrInternalError
	}

	logger.Info(ctx, "updated user", "user_id", userId)

	// Шаг 4. Если оффер бесплатный, сделать доставку того, что дает оффер
	if offer.IsFree {
		err = s.enrollUser(ctx, userId, dto.Email, offer)
		if err != nil {
			return nil, err
		}

		result.Message = *offer.SuccessMessage
		result.RedirectURL = *offer.RedirectURL

		logger.Info(ctx, "enrolled user", "user_id", userId, "offer_id", offer.ID)

		return &result, nil
	}

	// Шаг 4. Если оффер платный, создать заказ и вернуть ссылку на оплату
	payment, err := s.createPayment(ctx, CreatePaymentDTO{
		PayMethod:        offer.PayMethod,
		UserID:           dto.UserID,
		Email:            dto.Email,
		Phone:            dto.Phone,
		OrderDescription: offer.Name,
		OfferID:          offer.ID,
		Price:            offer.Price,
		// TODO: отправка письма о создании заказа может быть отключена
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	logger.Log.InfoContext(ctx, "created payment", "payment_id", payment.PaymentID, "order_id", payment.OrderID, "payment_url", payment.PaymentURL)

	result.RedirectURL = payment.PaymentURL

	return &result, nil
}

func (s *Service) createPayment(ctx context.Context, dto CreatePaymentDTO) (*payments.GetPaymentLinkResult, error) {
	if dto.PayMethod == nil {
		return nil, common.ErrInternalError
	}

	// создать заказ
	orderId, err := s.createOrder(ctx, dto.UserID, dto.PayMethod.ID, dto.OfferID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	// создать ссылку на оплату
	paymentSystem := payments.NewPaymentSystem(dto.PayMethod.Type)
	if paymentSystem == nil {
		return nil, common.ErrPaymentSystemNotFound
	}

	paymentResult, err := paymentSystem.GetPaymentLink(ctx, payments.GetPaymentLinkPayload{
		Email:           dto.Email,
		Description:     dto.OrderDescription,
		Amount:          dto.Price,
		Phone:           dto.Phone,
		Login:           dto.PayMethod.Login,
		Password:        dto.PayMethod.Password,
		OrderId:         orderId,
		SendReceipt:     dto.PayMethod.SendReceipt,
		ReceiptSettings: dto.PayMethod.ReceiptSettings,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	// Обновить paymentId в созданного заказа
	// paymentId у платежных систем генерируется со ссылкой для оплаты
	err = s.repo.UpdateOrderPaymentId(ctx, paymentResult.OrderID, paymentResult.PaymentID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	// отправить письмо, что заказ создан
	err = s.sendOrderCreatedEmail(ctx, dto.Email, dto.OrderDescription, dto.Price, paymentResult.PaymentURL)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	return paymentResult, nil
}

func (s *Service) createOrder(ctx context.Context, userId int64, payMethodId int64, offerId int64) (int64, error) {
	newOrder := NewOrder{
		UserID:        userId,
		OfferID:       offerId,
		IntegrationID: payMethodId,
	}

	orderId, err := s.repo.CreateOrder(ctx, newOrder)
	if err != nil {
		logger.Log.Error(err.Error())
		return 0, common.ErrInternalError
	}

	return orderId, nil
}

func (s *Service) enrollUser(ctx context.Context, userId int64, userEmail string, offer *OfferForProcessing) error {
	groups, err := s.repo.GetOfferGroups(ctx, offer.ID)
	if err != nil {
		logger.Log.ErrorContext(ctx, err.Error())
		return common.ErrInternalError
	}

	err = s.repo.AddUserToGroups(ctx, userId, groups)
	if err != nil {
		logger.Log.ErrorContext(ctx, err.Error())
		return common.ErrInternalError
	}

	if offer.SendRegistrationEmail {
		err = s.sendEnrollmentEmail(ctx, userEmail, *offer.RegistrationEmailTheme, *offer.RegistrationEmail)
		if err != nil {
			logger.Error(ctx, err.Error())
		}
	}

	return nil
}

func (s *Service) GetProfile(ctx context.Context, userId int) (*Profile, error) {
	profile, err := s.repo.GetProfileByUserId(ctx, userId)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	return profile, nil
}

func (s *Service) ProcessTinkoffWebhook(ctx context.Context, payload TinkoffWebhookBody) error {
	// Отформатировать статус
	status := payments.FormatStatus(payload.Status)

	orderId, err := strconv.ParseInt(payload.OrderId, 10, 64)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("could not parse int64 from orderId %s", payload.OrderId), "err", err.Error())
		return common.ErrInternalError
	}

	// Получить заказ
	order, err := s.repo.FindOrderById(ctx, orderId)
	if err != nil {
		if !errors.Is(err, common.ErrOrderNotFound) {
			logger.Error(ctx, fmt.Sprintf("could not get order by id", payload.OrderId), "err", err.Error())
			return common.ErrInternalError
		}
	}

	// Провалидировать данные
	err = s.validateTinkoffWebhook(ctx, payload, order)
	if err != nil {
		return err
	}

	logger.Info(ctx, "got valid order webhook", "orderId", orderId, "status", status)

	// Обновить заказ
	cardInfo := OrderCardInfo{
		ExpirationDate: payload.ExpDate,
		Pan:            payload.Pan,
	}

	orderError := OrderError{
		Message:    payload.Message,
		Details:    payload.Details,
		StatusCode: payload.ErrorCode,
	}

	if payload.ErrorCode == "" {
		orderError.StatusCode = "0"
	}

	err = s.repo.UpdateOrderStatus(ctx, order.ID, status, orderError, cardInfo)
	if err != nil {
		return common.ErrInternalError
	}

	if status == payments.StatusSucceeded {
		err = s.processSucceededOrder(ctx, order)
		return err
	}

	return nil
}

func (s *Service) processSucceededOrder(ctx context.Context, order *OrderForProcessing) error {

	// Выдать пользователю оффер
	offer, err := s.GetOfferForProcessing(ctx, order.OfferSlug)
	if err != nil {
		logger.Error(ctx, "could not find offer for processing", "order_id", order.ID, "offer_slug", order.OfferSlug, "err", err.Error())
		return common.ErrInternalError
	}

	err = s.sendOrderCompletedEmail(ctx, order.UserEmail, offer.Name, offer.Price)
	if err != nil {
		logger.Error(ctx, "could not send order completed email", "order", order.ID, "err", err.Error())
		return common.ErrInternalError
	}

	err = s.enrollUser(ctx, order.UserID, order.UserEmail, offer)
	if err != nil {
		logger.Error(ctx, "could not enroll user", "order", order.ID, "err", err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) validateTinkoffWebhook(ctx context.Context, payload TinkoffWebhookBody, order *OrderForProcessing) error {
	if order.PaymentID != strconv.Itoa(int(payload.PaymentId)) {
		logger.Error(ctx, "order payment id not equal with webhook payment id", "order_payment_id", order.PaymentID, "webhook_payment_id", payload.PaymentId)
		return common.ErrInternalError
	}

	if order.Price != payload.Amount {
		logger.Error(ctx, "order price not equal with webhook amount", "order_price", order.Price, "webhook_amount", payload.Amount, "order_id", order.ID)
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, userId int, profile UpdateProfileBody) error {
	err := s.repo.UpdateProfile(ctx, userId, profile)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) ChangeAvatar(ctx context.Context, userId int, avatarPathToDir string, avatarFileName string) error {
	defer func() {
		err := RemoveLocalFile(avatarPathToDir + "/" + avatarFileName)
		if err != nil {
			logger.Log.Error(err.Error(), "error", err)
		}
	}()

	src, err := imaging.Open(avatarPathToDir + "/" + avatarFileName)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	fileName := fmt.Sprintf("avatar_for_user_%d", userId)
	newAvatarFileName := MakeFileHashName(fileName, "jpeg")

	buff := new(bytes.Buffer)

	size := src.Bounds().Size()

	if size.X > 300 {
		src = imaging.Resize(src, 300, 0, imaging.Lanczos)
	}

	err = jpeg.Encode(buff, src, nil)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	fileUrl, err := UploadFileToS3(s.config.PhotosBucket, newAvatarFileName, bytes.NewReader(buff.Bytes()), s.config)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	err = s.repo.UpdateAvatar(ctx, userId, fileUrl)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userId int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	err = s.repo.UpdatePassword(ctx, userId, string(hashedPassword))
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) GetUserAccessibleProducts(ctx context.Context, userId int) ([]ProductCard, error) {
	products, err := s.repo.GetUserAccessibleProducts(ctx, userId)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	return products, nil
}

func (s *Service) GetUserAccessibleProduct(ctx context.Context, courseSlug string, userId int) (*ProductInfo, error) {
	product, err := s.repo.GetUserAccessibleProduct(ctx, courseSlug, userId)

	if errors.Is(err, common.ErrProductNotFound) {
		return nil, common.ErrProductNotFound
	}

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	if product == nil {
		return nil, common.ErrProductNotFound
	}

	lessons, err := s.repo.GetProductLessons(ctx, product.ID)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	product.Lessons = lessons

	return product, nil
}

func (s *Service) GetUserAccessibleLesson(ctx context.Context, lessonSlug string, userId int) (*LessonInfo, error) {
	lesson, err := s.repo.GetUserAccessibleLesson(ctx, lessonSlug, userId)

	if errors.Is(err, common.ErrLessonNotFound) {
		return nil, common.ErrLessonNotFound
	}

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	if lesson == nil {
		return nil, common.ErrLessonNotFound
	}

	return lesson, nil
}

func (s *Service) CompleteLesson(ctx context.Context, lessonSlug string, userId int) error {
	err := s.repo.CompleteLesson(ctx, lessonSlug, userId)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}
	return nil
}

func (s *Service) GetSolvedQuizzesForQuiz(ctx context.Context, lessonSlug string, skip int, limit int) ([]QuizSolvedInfo, error) {
	solvedQuizzes, err := s.repo.GetSolvedQuizzesForQuiz(ctx, lessonSlug, skip, limit)
	if err != nil {
		logger.Log.Error(err.Error())
		return solvedQuizzes, common.ErrInternalError
	}
	return solvedQuizzes, nil
}

func (s *Service) GetSolvedQuizzesForProduct(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error) {
	solvedQuizzes, err := s.repo.GetSolvedQuizzesForProduct(ctx, productSlug, userId, skip, limit)
	if err != nil {
		logger.Log.Error(err.Error())
		return solvedQuizzes, common.ErrInternalError
	}
	return solvedQuizzes, nil
}

func (s *Service) GetSolvedQuizzesForUser(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error) {
	solvedQuizzes, err := s.repo.GetSolvedQuizzesForUser(ctx, productSlug, userId, skip, limit)
	if err != nil {
		logger.Log.Error(err.Error())
		return solvedQuizzes, common.ErrInternalError
	}
	return solvedQuizzes, nil
}

func (s *Service) SolveQuiz(ctx context.Context, dto SolveQuizDTO) error {

	solvedQuiz, err := s.repo.FindSolvedQuiz(ctx, dto.UserID, dto.QuizSlug)

	if err != nil && !errors.Is(err, common.ErrSolvedQuizNotFound) {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	if solvedQuiz != nil {
		return common.ErrQuizAlreadySolved
	}

	var savedMediaIds []int64

	// Загружаем медиа у выполненного квиза, если они есть
	for _, file := range dto.Media {
		if file.MediaType == "image" {
			res, err := s.createImageMediaFromLocalFile(ctx, file)
			if err != nil {
				logger.Log.Error(err.Error())
				continue
			}
			savedMediaIds = append(savedMediaIds, res.MediaId)
		} else if file.MediaType == "video" {
			res, err := s.createVideoMediaFromLocalFile(ctx, file)
			if err != nil {
				logger.Log.Error(err.Error())
				continue
			}
			savedMediaIds = append(savedMediaIds, res.MediaId)
		} else {
			logger.Log.Error("unknown media type", "mediaType", file.MediaType)
		}
	}

	// Сохраняем выполненный квиз
	answerJson, err := json.Marshal(QuizSolvedAnswer{Answer: dto.Answer})
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	solvedQuizId, err := s.repo.SolveQuiz(ctx, dto.QuizSlug, dto.UserID, answerJson)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	// Привязываем медиа к выполненному квизу
	if len(savedMediaIds) > 0 {
		err = s.repo.ConnectManyMedia(ctx, savedMediaIds, RelatedMediaTypeSolvedQuiz, solvedQuizId)
		if err != nil {
			logger.Log.Error(err.Error(), "mediaIds", savedMediaIds, "solvedQuizId", solvedQuizId)
			return common.ErrInternalError
		}
	}

	return nil
}

func (s *Service) DeleteSolvedQuiz(ctx context.Context, quizSlug string, userId int) error {
	solvedQuiz, err := s.repo.FindSolvedQuiz(ctx, userId, quizSlug)
	if err != nil {
		if errors.Is(err, common.ErrSolvedQuizNotFound) {
			return nil
		}
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}
	err = s.repo.DeleteSolvedQuiz(ctx, solvedQuiz.ID, userId)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}
	return nil
}

func (s *Service) createVideoMediaFromLocalFile(ctx context.Context, file FileUpload) (FileUploadResult, error) {
	var result FileUploadResult

	f, err := os.Open(file.Path)
	if err != nil {
		logger.Log.Error(err.Error())
		return result, common.ErrInternalError
	}

	slug := uuid.New().String()
	ext := GetExtensionFromFileName(file.FileName)

	fileName := slug + "." + ext
	bucket := s.config.VideosBucket

	fileUrl, err := UploadFileToS3(bucket, fileName, f, s.config)

	if err != nil {
		logger.Log.Error(err.Error())
		return result, common.ErrInternalError
	}

	media := Media{
		Slug:    slug,
		URL:     fileUrl,
		Size:    &file.Size,
		Name:    fileName,
		Ext:     ext,
		Mime:    file.Mime,
		Bucket:  bucket,
		Storage: s.config.S3Provider,
		Type:    file.MediaType,
		Status:  MediaStatusUploaded,
	}

	mediaId, err := s.repo.SaveMedia(ctx, media)

	if err != nil {
		logger.Log.Error(err.Error(), "fileUrl", fileUrl)
		return result, common.ErrInternalError
	}

	logger.Log.Info("uploaded file url", "fileUrl", fileUrl)

	result.FileURL = fileUrl
	result.MediaId = mediaId

	return result, nil
}

func (s *Service) createImageMediaFromLocalFile(ctx context.Context, file FileUpload) (FileUploadResult, error) {
	src, err := imaging.Open(file.Path)

	var result FileUploadResult

	if err != nil {
		logger.Log.Error(err.Error())
		return result, common.ErrInternalError
	}

	slug := uuid.New().String()
	ext := "jpeg"
	mime := "image/jpeg"

	fileName := slug + "." + ext

	buff := new(bytes.Buffer)

	// Конвертация изображения в jpeg
	err = jpeg.Encode(buff, src, nil)
	if err != nil {
		logger.Log.Error(err.Error())
		return result, common.ErrInternalError
	}

	// Загрузка в S3
	bucket := s.config.PhotosBucket
	fileUrl, err := UploadFileToS3(bucket, fileName, bytes.NewReader(buff.Bytes()), s.config)

	if err != nil {
		logger.Log.Error(err.Error())
		return result, common.ErrInternalError
	}

	// Получение размеров изображения
	dimensions := src.Bounds().Size()

	// Сохранение медиа в БД
	media := Media{
		Slug:    slug,
		URL:     fileUrl,
		Width:   &dimensions.X,
		Height:  &dimensions.Y,
		Size:    &file.Size,
		Name:    fileName,
		Ext:     ext,
		Mime:    mime,
		Bucket:  bucket,
		Storage: s.config.S3Provider,
		Type:    file.MediaType,
		Status:  MediaStatusUploaded,
	}

	mediaId, err := s.repo.SaveMedia(ctx, media)

	if err != nil {
		logger.Log.Error(err.Error(), "fileUrl", fileUrl)
		return result, common.ErrInternalError
	}

	logger.Log.Info("uploaded file url", "fileUrl", fileUrl)

	result.FileURL = fileUrl
	result.MediaId = mediaId

	return result, nil
}

func (s *Service) GetQuizBySlug(ctx context.Context, slug string) (*Quiz, error) {
	quiz, err := s.repo.GetQuizBySlug(ctx, slug)

	if errors.Is(err, common.ErrQuizNotFound) {
		return nil, common.ErrQuizNotFound
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	return quiz, nil
}

func (s *Service) generatePassword() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?,."
	const passwordLength = 8

	charsetLength := big.NewInt(int64(len(chars)))
	password := make([]byte, passwordLength)

	for i := range password {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password[i] = chars[num.Int64()]
	}

	return string(password), nil
}

func (s *Service) createUserPassword(ctx context.Context, userPassword string) (string, string, error) {
	// сгенерировать пароль — если отсутствует
	if userPassword == "" {
		password, err := s.generatePassword()
		if err != nil {
			logger.Log.Error(err.Error())
			return "", "", nil
		}
		userPassword = password
	}

	// заэшировать пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), 10)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", "", nil
	}

	return string(hashedPassword), userPassword, nil
}

func (s *Service) createUser(ctx context.Context, dto CreateUserDTO) (int64, bool, error) {

	// создать пароль для пользователя
	hashedPassword, rawPassword, err := s.createUserPassword(ctx, dto.Password)

	if err != nil {
		logger.Log.Error(err.Error())
		return 0, false, common.ErrInternalError
	}

	// создать пользователя в базе
	user := User{
		Email:     dto.Email,
		Password:  hashedPassword,
		FirstName: dto.FirstName,
	}

	userId, err := s.repo.CreateUser(ctx, user)

	var alreadyExists bool

	if err != nil {
		if !errors.Is(err, common.ErrUserAlreadyExists) {
			logger.Log.Error(err.Error(), "error", err)
			return userId, alreadyExists, err
		}

		alreadyExists = true
	}

	if alreadyExists {
		return userId, alreadyExists, nil
	}

	// отправить welcome-письмо
	err = s.sendWelcomeEmail(ctx, user.Email, rawPassword)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	return userId, alreadyExists, err
}

func (s *Service) sendOrderCompletedEmail(ctx context.Context, userEmail string, ordered string, amount uint64) error {
	email, err := s.emails.GetEmailByType(ctx, "order-completed")

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
	}

	email.Context["Ordered"] = ordered
	email.Context["Amount"] = amount
	email.Context["HeroURL"] = s.config.HeroAppBaseURL + "/login?way=password&email=" + userEmail

	err = s.emails.SendEmail(email, []string{userEmail})

	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) sendWelcomeEmail(ctx context.Context, userEmail string, userPassword string) error {
	email, err := s.emails.GetEmailByType(ctx, "welcome")

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
	}

	email.Context["Email"] = userEmail
	email.Context["Password"] = userPassword
	email.Context["LoginURL"] = s.config.HeroAppBaseURL + "/login"
	email.Context["LoginFullURL"] = s.config.HeroAppBaseURL + "/login?way=password&email=" + userEmail
	email.Context["MailFrom"] = email.From.Email

	err = s.emails.SendEmail(email, []string{userEmail})

	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) sendOrderCreatedEmail(ctx context.Context, userEmail string, ordered string, amount uint64, paymentUrl string) error {
	email, err := s.emails.GetEmailByType(ctx, "order-created")

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
	}

	email.Context["PaymentURL"] = paymentUrl
	email.Context["Ordered"] = ordered
	email.Context["Amount"] = amount

	err = s.emails.SendEmail(email, []string{userEmail})

	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) sendEnrollmentEmail(ctx context.Context, userEmail string, emailSubject string, emailBody string) error {
	email, err := s.emails.GetEmailByType(ctx, "general")

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	email.Subject = emailSubject
	email.Body = emailBody

	err = s.emails.SendEmail(email, []string{userEmail})

	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	logger.Log.InfoContext(ctx, "send enrollment email", "user_email", userEmail, "email_subject", emailSubject)

	return nil
}

func (s *Service) Signup(ctx context.Context, body *SignupBody) (*SignUpResult, error) {

	_, alreadyExists, err := s.createUser(ctx, CreateUserDTO{
		FirstName: body.FirstName,
		Email:     body.Email,
		Password:  body.Password,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}

	// Находим только что созданного пользователя или уже существующего
	foundUser, err := s.repo.FindUserByEmail(ctx, body.Email)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	var result SignUpResult

	if !alreadyExists {
		result.Message = "Регистрация прошла успешно! На твою почту было отправлено письмо с данными для входа"
		return &result, nil
	}

	// Если пользователь уже существовал и его пароль при попытке создать аккаунт
	// Совпадает с тем, который есть в базе — значит, можем его авторизовать
	// А если не совпадает, то отдает результат — что такой аккаунт уже есть
	if !s.passwordMatches(foundUser.Password, body.Password) {
		result.Message = "Оу, оказывается, у тебя уже есть аккаунт"
		return &result, nil
	}

	token, err := s.createJWTToken(foundUser.ID)

	if err != nil {
		return &result, common.ErrInternalError
	}

	result.Token = &token
	result.Message = "Привет. С возвращением!"

	return &result, nil

}

func (s *Service) Login(ctx context.Context, body *LoginBody) (*LoginResult, error) {

	user, err := s.repo.FindUserByEmail(ctx, body.Email)

	if err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return nil, common.ErrWrongCredentials
		}
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	if !s.passwordMatches(user.Password, body.Password) {
		return nil, common.ErrWrongCredentials
	}

	token, err := s.createJWTToken(user.ID)

	if err != nil {
		return nil, common.ErrInternalError
	}
	result := LoginResult{
		Token: token,
	}

	return &result, nil

}

func (s *Service) GetMagicLink(ctx context.Context, to string) error {
	user, err := s.repo.FindUserByEmail(ctx, to)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}
	if user == nil {
		return common.ErrUserNotFound
	}

	magicLink, err := s.createMagicLink(user.ID)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.ErrInternalError
	}

	type MagicLinkContent struct {
		MagicLink string `json:"magic_link"`
	}

	email, err := s.emails.GetEmailByType(ctx, "magiclink")

	if err != nil {
		return err
	}

	email.Context["MagicLink"] = magicLink

	err = s.emails.SendEmail(email, []string{to})

	if err != nil {
		logger.Log.Error(err.Error())
		return common.ErrInternalError
	}

	return nil
}

func (s *Service) ValidateMagicLink(ctx context.Context, token string) (*LoginResult, error) {
	user, err := s.ValidateJWTToken(ctx, token)
	if err != nil {
		if errors.Is(err, common.ErrTokenExpired) {
			return nil, common.ErrMagicLinkExpired
		}
		if errors.Is(err, common.ErrInvalidToken) {
			return nil, common.ErrInvalidMagicLink
		}
		return nil, err
	}

	jwtToken, err := s.createJWTToken(user.ID)

	if err != nil {
		return nil, common.ErrInternalError
	}

	result := LoginResult{
		Token: jwtToken,
	}

	return &result, nil

}

func (s *Service) createMagicLink(userId int) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.MagicLinkExp)),
		},
		UserID: userId,
	}

	token := jwt.NewWithClaims(s.config.JwtSigningMethod, claims)

	tokenString, err := token.SignedString([]byte(s.config.JwtTokenSecretKey))

	if err != nil {
		return "", err
	}

	magicLink := s.config.HeroAppBaseURL + "/login/magic-link?token=" + tokenString

	return magicLink, nil
}

func (s *Service) createJWTToken(userId int) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JwtTokenExp)),
		},
		UserID: userId,
	}

	token := jwt.NewWithClaims(s.config.JwtSigningMethod, claims)

	tokenString, err := token.SignedString([]byte(s.config.JwtTokenSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateJWTToken(ctx context.Context, token string) (*User, error) {
	claims := Claims{}
	data, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.config.JwtTokenSecretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, common.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, common.ErrInvalidToken
		}
		return nil, err
	}

	if !data.Valid {
		return nil, common.ErrInvalidToken
	}

	if data.Method != s.config.JwtSigningMethod {
		logger.Log.Warn(fmt.Sprintf("JWT Token method mismatch"))
		return nil, common.ErrInvalidToken
	}

	user, err := s.repo.FindUserById(ctx, claims.UserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) passwordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewService(repo Storage, config *config.Config, emails IEmailsService) *Service {
	return &Service{
		repo:   repo,
		config: config,
		emails: emails,
	}
}
