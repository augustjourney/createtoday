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

	// products
	GetUserAccessibleProducts(ctx context.Context, userId int) ([]ProductCard, error)
	GetUserAccessibleProduct(ctx context.Context, productSlug string, userId int) (*ProductInfo, error)
	GetProductLessons(ctx context.Context, productId int) ([]LessonCard, error)

	// lessons
	GetUserAccessibleLesson(ctx context.Context, lessonSlug string, userId int) (*LessonInfo, error)
}
