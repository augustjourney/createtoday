package hero

import (
	"context"
	"createtodayapi/internal/entity"
)

type Storage interface {
	// users
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	FindUserById(ctx context.Context, id int) (*entity.User, error)

	// profile
	GetProfileByUserId(ctx context.Context, userId int) (*entity.Profile, error)

	// products
	GetUserAccessibleProducts(ctx context.Context, userId int) ([]entity.UserProductCard, error)
}
