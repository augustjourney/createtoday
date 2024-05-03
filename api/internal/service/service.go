package service

import (
	"context"
	"createtodayapi/internal/dto"
	"createtodayapi/internal/entity"
)

type Auth interface {
	Signup(ctx context.Context, body *dto.SignupBody) (*dto.SignUpResult, error)
	Login(ctx context.Context, body *dto.LoginBody) (*dto.LoginResult, error)
	GetMagicLink(ctx context.Context, to string) error
	ValidateMagicLink(ctx context.Context, token string) (*dto.LoginResult, error)
	ValidateJWTToken(ctx context.Context, token string) (*entity.User, error)

	createMagicLink(userId int) (string, error)
	createJWTToken(userId int) (string, error)
	passwordMatches(hash string, password string) bool
}

type Emails interface {
	GetEmailByType(context context.Context, emailType string) (*entity.Email, error)
	SendEmail(email *entity.Email, to []string) error

	buildTemplatePath(templateName string) (string, error)
	buildEmailHtml(email *entity.Email) (string, error)
	buildEmailSenderName(sender entity.EmailSender) string
}

type Profile interface {
	GetProfile(ctx context.Context, userId int) (*entity.Profile, error)
}
