package hero

import (
	"context"
)

type Storage interface {
	// users
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user User) error
	FindUserById(ctx context.Context, id int) (*User, error)

	// profile
	GetProfileByUserId(ctx context.Context, userId int) (*Profile, error)
	UpdateProfile(ctx context.Context, userId int, profile UpdateProfileBody) error
	UpdateAvatar(ctx context.Context, userId int, avatar string) error
	UpdatePassword(ctx context.Context, userId int, password string) error

	// products
	GetUserAccessibleProducts(ctx context.Context, userId int) ([]ProductCard, error)
	GetUserAccessibleProduct(ctx context.Context, productSlug string, userId int) (*ProductInfo, error)
	GetProductLessons(ctx context.Context, productId int) ([]LessonCard, error)

	// lessons
	GetUserAccessibleLesson(ctx context.Context, lessonSlug string, userId int) (*LessonInfo, error)

	// quizzes
	GetSolvedQuizzesForQuiz(ctx context.Context, quizSlug string, skip int, limit int) ([]QuizSolvedInfo, error)
	GetSolvedQuizzesForProduct(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error)
	GetSolvedQuizzesForUser(ctx context.Context, productSlug string, userId int, skip int, limit int) ([]QuizSolvedInfo, error)
	SolveQuiz(ctx context.Context, quizSlug string, userId int, answer []byte) (int64, error)
	FindSolvedQuiz(ctx context.Context, userId int, quizSlug string) (*QuizSolved, error)
	GetQuizBySlug(ctx context.Context, quizSlug string) (*Quiz, error)
	DeleteSolvedQuiz(ctx context.Context, solvedQuizId int64, userId int) error

	// media
	SaveMedia(ctx context.Context, media Media) (int64, error)
	ConnectMedia(ctx context.Context, mediaId int64, relatedType string, relatedId int64) error
	ConnectManyMedia(ctx context.Context, mediaIds []int64, relatedType string, relatedId int64) error
	DeleteMedia(ctx context.Context, mediaId int64) error
	UpdateMediaStatus(ctx context.Context, mediaId int64, status string) error
}
