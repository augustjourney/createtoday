package dto

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

type LoginResult struct {
	Token string `json:"token"`
}

type SignUpResult struct {
	AlreadyExists bool   `json:"alreadyExists"`
	Token         string `json:"token"`
}
