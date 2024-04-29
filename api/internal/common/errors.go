package common

import "errors"

var ErrInternalError = errors.New("Что-то пошло не так")
var ErrUserAlreadyExists = errors.New("Такой пользователь уже существует")
var ErrUserNotFound = errors.New("Пользователь не найден")
var ErrWrongCredentials = errors.New("Неверный пароль или логин")
var ErrEmptyEmail = errors.New("Email не может быть пустым")
var ErrEmptyPassword = errors.New("Пароль не может быть пустым")
