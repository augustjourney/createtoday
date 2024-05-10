package hero

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type IController interface {
	// Auth
	Login(ctx *fiber.Ctx) error
	Signup(ctx *fiber.Ctx) error
	GetMagicLink(ctx *fiber.Ctx) error
	ValidateMagicLink(ctx *fiber.Ctx) error

	// Products
	GetUserAccessibleProducts(ctx *fiber.Ctx) error

	// Profile
	GetProfile(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error

	// Quizzes
	SolveQuiz(ctx *fiber.Ctx) error
}

type Controller struct {
	service IService
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "application/json")
	var body LoginBody

	err := json.Unmarshal(ctx.Body(), &body)

	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = body.Validate()

	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	result, err := c.service.Login(context.Background(), &body)

	if errors.Is(err, common.ErrWrongCredentials) {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, common.ErrInternalError)
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = result.Token

	ctx.Cookie(tokenCookie)

	return common.DoApiResponse(ctx, http.StatusOK, result, nil)
}

func (c *Controller) Signup(ctx *fiber.Ctx) error {
	var body SignupBody

	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	result, err := c.service.Signup(context.Background(), &body)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, common.ErrInternalError)
	}

	if result.Token != nil {
		tokenCookie := new(fiber.Cookie)
		tokenCookie.Name = "token"
		tokenCookie.Value = *result.Token

		ctx.Cookie(tokenCookie)
	}

	return common.DoApiResponse(ctx, http.StatusOK, result, nil)
}

func (c *Controller) GetMagicLink(ctx *fiber.Ctx) error {
	var body GetMagicLinkBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = c.service.GetMagicLink(context.Background(), body.Email)

	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	type GetMagicLinkResult struct {
		Message string `json:"message"`
	}

	return common.DoApiResponse(ctx, http.StatusOK, GetMagicLinkResult{
		Message: "Письмо с ссылкой для входа без пароля отправлено на вашу почту",
	}, nil)
}

func (c *Controller) ValidateMagicLink(ctx *fiber.Ctx) error {
	var body ValidateMagicLinkBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	result, err := c.service.ValidateMagicLink(context.Background(), body.Token)

	if err != nil {
		logger.Log.Error(err.Error())
		if errors.Is(err, common.ErrMagicLinkExpired) || errors.Is(err, common.ErrInvalidMagicLink) {
			return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
		}
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = result.Token

	ctx.Cookie(tokenCookie)

	return common.DoApiResponse(ctx, http.StatusOK, result, nil)
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	profile, err := c.service.GetProfile(context.Background(), user.ID)

	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	if profile == nil {
		return common.DoApiResponse(ctx, http.StatusNotFound, nil, common.ErrUserNotFound)
	}

	return common.DoApiResponse(ctx, http.StatusOK, profile, nil)
}

func (c *Controller) UpdateProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	var body UpdateProfileBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	// TODO: валидация body
	// Все поля из структуры должны быть в json-body
	// Например: если в json-body нет полей,
	// Которые должны быть в структуре
	// А в БД эти поля заполнены
	// Получается перезатрем их на null
	// При этом если поле есть в json-body со значением null
	// Это корректно

	err = c.service.UpdateProfile(context.Background(), user.ID, body)

	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	return common.DoApiResponse(ctx, http.StatusOK, "Профиль обновлен", nil)
}

func (c *Controller) ChangeAvatar(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)

	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	file := form.File["avatar"][0]

	if file == nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, common.ErrEmptyAvatar)
	}

	wd, err := os.Getwd()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, common.ErrInternalError)
	}

	avatarFileName := fmt.Sprintf("avatar_%d_%s", user.ID, file.Filename)
	avatarPathToDir := fmt.Sprintf("%s/temp", wd)

	err = ctx.SaveFile(file, avatarPathToDir+"/"+avatarFileName)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	err = c.service.ChangeAvatar(context.Background(), user.ID, avatarPathToDir, avatarFileName)

	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	return common.DoApiResponse(ctx, http.StatusOK, "Аватар успешно загружен", nil)

}

func (c *Controller) UpdatePassword(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	var body UpdatePasswordBody

	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	err = c.service.ChangePassword(context.Background(), user.ID, body.Password)
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	return common.DoApiResponse(ctx, http.StatusOK, "Новый пароль успешно сохранен", nil)
}

func (c *Controller) GetUserAccessibleProducts(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	products, err := c.service.GetUserAccessibleProducts(context.Background(), user.ID)
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}
	return common.DoApiResponse(ctx, http.StatusOK, products, nil)
}

func (c *Controller) GetUserAccessibleProduct(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	slug := ctx.Params("slug")
	product, err := c.service.GetUserAccessibleProduct(context.Background(), slug, user.ID)
	if errors.Is(err, common.ErrProductNotFound) {
		return common.DoApiResponse(ctx, http.StatusNotFound, nil, common.ErrProductNotFound)
	}
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}
	return common.DoApiResponse(ctx, http.StatusOK, product, nil)
}

func (c *Controller) GetUserAccessibleLesson(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	slug := ctx.Params("slug")
	lesson, err := c.service.GetUserAccessibleLesson(context.Background(), slug, user.ID)
	if errors.Is(err, common.ErrLessonNotFound) {
		return common.DoApiResponse(ctx, http.StatusNotFound, nil, err)
	}
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}
	return common.DoApiResponse(ctx, http.StatusOK, lesson, nil)
}

func (c *Controller) GetSolvedQuizzesForQuiz(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	solvedQuizzes, err := c.service.GetSolvedQuizzesForQuiz(context.Background(), slug)
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	return common.DoApiResponse(ctx, http.StatusOK, solvedQuizzes, nil)
}

func (c *Controller) SolveQuiz(ctx *fiber.Ctx) error {
	var body SolveQuizBody

	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	body.Type = form.Value["type"][0]
	body.Answer = form.Value["answer"][0]

	media := form.File["media"]

	defer func() {
		if len(body.Media) == 0 {
			return
		}

		for _, m := range body.Media {
			_ = RemoveLocalFile(m.Path)
		}
	}()

	if media != nil {
		wd, err := os.Getwd()
		if err != nil {
			logger.Log.Error(err.Error())
			return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, common.ErrInternalError)
		}

		for _, file := range media {
			filePath := fmt.Sprintf("%s/temp/%s", wd, file.Filename)

			err = ctx.SaveFile(file, filePath)
			if err != nil {
				logger.Log.Error(err.Error())
				return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
			}

			body.Media = append(body.Media, FileUpload{
				FileName: file.Filename,
				Size:     file.Size,
				Path:     filePath,
				Mime:     file.Header.Get("Content-Type"),
			})
		}
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	user := ctx.Locals("user").(*User)
	slug := ctx.Params("slug")

	err = c.service.SolveQuiz(context.Background(), SolveQuizDTO{
		Answer:   body.Answer,
		UserID:   user.ID,
		Type:     body.Type,
		QuizSlug: slug,
		Media:    body.Media,
	})

	if errors.Is(err, common.ErrQuizAlreadySolved) {
		return common.DoApiResponse(ctx, http.StatusBadRequest, nil, err)
	}

	if err != nil {
		return common.DoApiResponse(ctx, http.StatusInternalServerError, nil, err)
	}

	defer func() {
		for _, file := range body.Media {
			err := RemoveLocalFile(file.Path)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}()

	return common.DoApiResponse(ctx, http.StatusOK, "Задание успешно выполнено", nil)
}

func NewController(service IService) *Controller {
	return &Controller{
		service: service,
	}
}
