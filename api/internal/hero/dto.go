package hero

import (
	"createtodayapi/internal/common"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *LoginBody) Validate() error {

	if b.Email == "" {
		return common.ErrEmptyEmail
	}

	if b.Password == "" {
		return common.ErrEmptyPassword
	}

	return nil
}

type GetMagicLinkBody struct {
	Email string `json:"email"`
}

func (b *GetMagicLinkBody) Validate() error {
	if b.Email == "" {
		return common.ErrEmptyEmail
	}
	return nil
}

type ValidateMagicLinkBody struct {
	Token string `json:"token"`
}

func (b *ValidateMagicLinkBody) Validate() error {
	if b.Token == "" {
		return common.ErrInvalidToken
	}

	return nil
}

type LoginResult struct {
	Token string `json:"token"`
}

type SignupBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
}

func (b *SignupBody) Validate() error {

	if b.Email == "" {
		return common.ErrEmptyEmail
	}

	return nil
}

type SignUpResult struct {
	AlreadyExists bool    `json:"alreadyExists"`
	Token         *string `json:"token"`
	Message       string  `json:"message,omitempty"`
}

type UpdateProfileBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Telegram  *string `json:"telegram"`
	Instagram *string `json:"instagram"`
	About     *string `json:"about"`
	Phone     *string `json:"phone"`
}

type UpdatePasswordBody struct {
	Password string `json:"password"`
}

func (b *UpdatePasswordBody) Validate() error {
	if b.Password == "" {
		return common.ErrNewPasswordIsEmpty
	}

	if len(b.Password) < 8 {
		return common.ErrNewPasswordIsShort
	}

	return nil
}
