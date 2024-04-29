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
}

func (s *Auth) Signup(ctx context.Context, email string, password string) (*dto.SignUpResult, error) {
	result := dto.SignUpResult{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return &result, common.ErrInternalError
	}

	user := entity.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, user)

	if err != nil {

		// Если ошибка — не является `пользователь уже существует`
		// То выходим
		if !errors.Is(err, common.ErrUserAlreadyExists) {
			logger.Log.Error(err.Error(), "error", err)
			return &result, common.ErrInternalError
		}

		result.AlreadyExists = true
	}

	// Находим только что созданного пользователя или уже существующего
	foundUser, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return &result, common.ErrInternalError
	}

	// Если пользователь уже существовал и его пароль при попытке создать аккаунт
	// Совпадает с тем, который есть в базе — значит, можем его авторизовать
	// А если не совпадает, то отдает результат — что такой аккаунт уже есть
	if result.AlreadyExists {
		if !s.passwordMatches(foundUser.Password, password) {
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

func (s *Auth) Login(ctx context.Context, email string, password string) (*dto.LoginResult, error) {
	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return nil, common.ErrInternalError
	}

	if user == nil {
		return nil, common.ErrUserNotFound
	}

	if !s.passwordMatches(user.Password, password) {
		return nil, common.ErrWrongCredentials
	}

	token, err := s.createJWTToken(user.ID)

	if err != nil {
		return nil, common.ErrInternalError
	}

	return &dto.LoginResult{Token: token}, nil

}

func (s *Auth) createJWTToken(userId int) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT_TOKEN_EXP)),
		},
		UserID: userId,
	}

	token := jwt.NewWithClaims(s.config.JWT_SIGNING_METHOD, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWT_TOKEN_SECRET_KEY))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Auth) ValidateJWTToken(ctx context.Context, token string) (error, *entity.User) {
	claims := Claims{}
	data, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.config.JWT_TOKEN_SECRET_KEY), nil
		},
	)

	if err != nil {
		return err, nil
	}

	if !data.Valid {
		return nil, nil
	}

	if data.Method != s.config.JWT_SIGNING_METHOD {
		logger.Log.Warn(fmt.Sprintf("JWT Token method mismatch"))
		return nil, nil
	}

	user, err := s.repo.FindById(ctx, claims.UserID)

	if err != nil {
		return err, nil
	}

	return nil, user
}

func (s *Auth) passwordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAuthService(repo storage.Users, config *config.Config) *Auth {
	return &Auth{
		repo:   repo,
		config: config,
	}
}
