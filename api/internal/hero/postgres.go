package hero

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/entity"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const UsersTable = "public.user"
const ProjectsTable = "public.project"
const GroupsTable = "public.group"
const ProductsTable = "public.product"
const UserGroupsTable = "public.user_group"
const ProductGroupsTable = "public.product_group"
const UsersProductsView = "public._userproducts"

type PostgresRepo struct {
	db *sqlx.DB
}

func (r *PostgresRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	q := fmt.Sprintf(`select id, email, password from %s where email = $1`, UsersTable)
	var user entity.User
	err := r.db.GetContext(ctx, &user, q, email)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	q := fmt.Sprintf(`select id, email from %s where id = $1`, UsersTable)
	var user entity.User
	err := r.db.GetContext(ctx, &user, q, id)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) GetProfileByUserId(ctx context.Context, userId int) (*entity.Profile, error) {
	q := fmt.Sprintf(`
		select email, first_name, last_name, phone, avatar, telegram, instagram 
		from %s where id = $1`,
		UsersTable,
	)
	var profile entity.Profile
	err := r.db.GetContext(ctx, &profile, q, userId)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user entity.User) error {
	q := fmt.Sprintf(`
		insert into %s (email, password, first_name)
		values ($1, $2, $3)
	`, UsersTable)
	_, err := r.db.ExecContext(ctx, q, user.Email, user.Password, user.FirstName)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return common.ErrUserAlreadyExists
			}
		}
		return err
	}

	return err
}

func (r *PostgresRepo) GetUserAccessibleProducts(ctx context.Context, userId int) ([]entity.UserProductCard, error) {
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

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}
