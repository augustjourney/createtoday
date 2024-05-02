package postgres

import (
	"context"
	"createtodayapi/internal/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProductsRepo struct {
	db *sqlx.DB
}

func (r *ProductsRepo) GetUsersProducts(ctx context.Context, userId int) ([]entity.UserProductCard, error) {
	q := fmt.Sprintf("select distinct on (id) name, description, slug, cover, settings from %s where user_id = $1", UsersProductsView)
	var products []entity.UserProductCard
	err := r.db.SelectContext(ctx, &products, q, userId)
	if err != nil {
		return products, err
	}
	if products == nil {
		return []entity.UserProductCard{}, nil
	}
	return products, nil
}

func NewProductsRepo(db *sqlx.DB) *ProductsRepo {
	return &ProductsRepo{
		db: db,
	}
}
