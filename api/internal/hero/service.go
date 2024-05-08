package hero

import (
	"bytes"
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/config"
	"createtodayapi/internal/logger"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"image/jpeg"
	"time"
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

	ChangeAvatar(ctx context.Context, userId int, avatarPath string, avatarFileName string) error
	ChangePassword(ctx context.Context, userId int, password string) error

	GetSolvedQuizzesForQuiz(ctx context.Context, lessonSlug string) ([]QuizSolvedInfo, error)
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type Service struct {
	repo   Storage
	config *config.Config
	emails IEmailsService
}

func (s *Service) GetProfile(ctx context.Context, userId int) (*Profile, error) {
	profile, err := s.repo.GetProfileByUserId(ctx, userId)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	return profile, nil
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

func (s *Service) GetSolvedQuizzesForQuiz(ctx context.Context, lessonSlug string) ([]QuizSolvedInfo, error) {
	solvedQuizzes, err := s.repo.GetSolvedQuizzesForQuiz(ctx, lessonSlug)
	if err != nil {
		logger.Log.Error(err.Error())
		return solvedQuizzes, common.ErrInternalError
	}
	return solvedQuizzes, nil
}

func (s *Service) Signup(ctx context.Context, body *SignupBody) (*SignUpResult, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		return nil, common.ErrInternalError
	}

	user := User{
		Email:     body.Email,
		Password:  string(hashedPassword),
		FirstName: body.FirstName,
	}

	result := SignUpResult{}

	err = s.repo.CreateUser(ctx, user)

	if err != nil {

		// Если ошибка — не является `пользователь уже существует`
		// То выходим
		if !errors.Is(err, common.ErrUserAlreadyExists) {
			logger.Log.Error(err.Error(), "error", err)
			return nil, common.ErrInternalError
		}

		result.AlreadyExists = true
	}

	// Находим только что созданного пользователя или уже существующего
	foundUser, err := s.repo.FindUserByEmail(ctx, body.Email)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return nil, common.ErrInternalError
	}

	if !result.AlreadyExists {
		email, err := s.emails.GetEmailByType(ctx, "welcome")

		if err != nil {
			logger.Log.Error(err.Error(), "error", err)
		}

		email.Context["Email"] = body.Email
		email.Context["Password"] = body.Password
		email.Context["LoginURL"] = s.config.HeroAppBaseURL + "/login"
		email.Context["LoginFullURL"] = s.config.HeroAppBaseURL + "/login?way=password&email=" + body.Email
		email.Context["MailFrom"] = email.From.Email

		err = s.emails.SendEmail(email, []string{body.Email})

		if err != nil {
			logger.Log.Error(err.Error())
			return nil, common.ErrInternalError
		}

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
		if errors.As(err, &jwt.ErrTokenMalformed) {
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
