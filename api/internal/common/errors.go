package common

import "errors"

// common
var ErrInternalError = errors.New("Что-то пошло не так")

// users
var ErrUserAlreadyExists = errors.New("Такой пользователь уже существует")
var ErrUserNotFound = errors.New("Пользователь не найден")

// auth
var ErrWrongCredentials = errors.New("Неверный пароль или логин")
var ErrEmptyEmail = errors.New("Email не может быть пустым")
var ErrEmptyPassword = errors.New("Пароль не может быть пустым")
var ErrInvalidToken = errors.New("Неверный токен")
var ErrTokenExpired = errors.New("Сессия истекла")
var ErrMagicLinkExpired = errors.New("Время действия ссылки вышло")
var ErrInvalidMagicLink = errors.New("Некорректная ссылка")

// projects
var ErrProjectAlreadyExists = errors.New("Такой проект уже существует")
var ErrProjectNotFound = errors.New("Проект не найден")

// products
var ErrProductNotFound = errors.New("Такой курс не найден или у вас нет к нему доступа")

// lessons
var ErrLessonNotFound = errors.New("Такой урок не найден или у вас нет к нему доступа")

// avatar
var ErrEmptyAvatar = errors.New("Аватар не может быть пустым")
