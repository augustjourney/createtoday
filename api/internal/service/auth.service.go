package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"createtodayapi/internal/common"
	"createtodayapi/internal/config"
	"createtodayapi/internal/dto"
	"createtodayapi/internal/entity"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/storage"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type Auth struct {
	config *config.Config
	repo   storage.Users
	emails *EmailsService
}

func (s *Auth) Signup(ctx context.Context, body *dto.SignupBody) (*dto.SignUpResult, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		return nil, common.ErrInternalError
	}

	user := entity.User{
		Email:     body.Email,
		Password:  string(hashedPassword),
		FirstName: body.FirstName,
	}

	result := dto.SignUpResult{}

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
	foundUser, err := s.repo.FindByEmail(ctx, body.Email)

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
	}

	// Если пользователь уже существовал и его пароль при попытке создать аккаунт
	// Совпадает с тем, который есть в базе — значит, можем его авторизовать
	// А если не совпадает, то отдает результат — что такой аккаунт уже есть
	if result.AlreadyExists {
		if !s.passwordMatches(foundUser.Password, body.Password) {
			return &result, nil
		}
	}

	token, err := s.createJWTToken(foundUser.ID)

	if err != nil {
		return &result, common.ErrInternalError
	}

	result.Token = token

	return &result, nil

}

func (s *Auth) Login(ctx context.Context, body *dto.LoginBody) (*dto.LoginResult, error) {

	user, err := s.repo.FindByEmail(ctx, body.Email)

	if err != nil {
		return nil, common.ErrInternalError
	}

	if user == nil {
		return nil, common.ErrUserNotFound
	}

	if !s.passwordMatches(user.Password, body.Password) {
		return nil, common.ErrWrongCredentials
	}

	token, err := s.createJWTToken(user.ID)

	if err != nil {
		return nil, common.ErrInternalError
	}
	result := dto.LoginResult{
		Token: token,
	}

	return &result, nil

}

func (s *Auth) GetMagicLink(ctx context.Context, to string) error {
	user, err := s.repo.FindByEmail(ctx, to)
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

func (s *Auth) ValidateMagicLink(ctx context.Context, token string) (*dto.LoginResult, error) {
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

	result := dto.LoginResult{
		Token: jwtToken,
	}

	return &result, nil

}

func (s *Auth) createMagicLink(userId int) (string, error) {

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

func (s *Auth) createJWTToken(userId int) (string, error) {

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

func (s *Auth) ValidateJWTToken(ctx context.Context, token string) (*entity.User, error) {
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

	user, err := s.repo.FindById(ctx, claims.UserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Auth) passwordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAuthService(repo storage.Users, config *config.Config, emailService *EmailsService) *Auth {
	return &Auth{
		repo:   repo,
		config: config,
		emails: emailService,
	}
}
