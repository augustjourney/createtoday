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

// profile
var ErrNewPasswordIsEmpty = errors.New("Новый пароль не может быть пустым")
var ErrNewPasswordIsShort = errors.New("Новый пароль не может быть меньше 8 символов")

// quizzes
var ErrEmptyQuizType = errors.New("Тип задания не может быть пустым")
var ErrEmptyQuizAnswer = errors.New("Ответ не может быть пустым")
var ErrEmptyQuizPhoto = errors.New("Чтобы выполнить задание, загрузите фото")
var ErrEmptyQuizVideo = errors.New("Чтобы выполнить задание, загрузите видео")
var ErrQuizTooManyPhotos = errors.New("Загрузить можно только одно фото к этому заданию")
var ErrQuizTooManyVideos = errors.New("Загрузить можно только одно видео к этому заданию")
var ErrQuizAlreadySolved = errors.New("Вы уже выполнили это задание")
var ErrSolvedQuizNotFound = errors.New("Не найден такой выполненный квиз")
var ErrQuizNotFound = errors.New("Квиз не найден")

// offers
var ErrOfferNotFound = errors.New("Такой оффер не найден")

// payments
var ErrPaymentSystemNotFound = errors.New("Такой платежный метод не найден")

// orders
var ErrOrderNotFound = errors.New("Такой заказ не найден")

// cache
var ErrCacheItemNotFound = errors.New("Такое ключ не найден в кэше")
