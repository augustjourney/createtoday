package common

import "errors"

var ErrInternalError = errors.New("Что-то пошло не так")
var ErrUserAlreadyExists = errors.New("Такой пользователь уже существует")
var ErrUserNotFound = errors.New("Пользователь не найден")
var ErrWrongCredentials = errors.New("Неверный пароль или логин")
var ErrEmptyEmail = errors.New("Email не может быть пустым")
var ErrEmptyPassword = errors.New("Пароль не может быть пустым")
var ErrInvalidToken = errors.New("Неверный токен")
var ErrTokenExpired = errors.New("Сессия истекла")
var ErrMagicLinkExpired = errors.New("Время действия ссылки вышло")
var ErrInvalidMagicLink = errors.New("Некорректная ссылка")
