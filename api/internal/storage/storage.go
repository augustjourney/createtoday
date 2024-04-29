package storage

import (
	"context"
	"createtodayapi/internal/entity"
)

type Users interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	FindById(ctx context.Context, id int) (*entity.User, error)
	GetProfileByUserId(ctx context.Context, userId int) (*entity.Profile, error)
}
