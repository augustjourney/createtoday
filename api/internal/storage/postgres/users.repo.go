package postgres

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/entity"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UsersRepo struct {
	db *sqlx.DB
}

const UsersTable = "public.user"

func (r *UsersRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	q := fmt.Sprintf(`select id, email, password from %s where email = $1`, UsersTable)
	var user entity.User
	err := r.db.Get(&user, q, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) FindById(ctx context.Context, id int) (*entity.User, error) {
	q := fmt.Sprintf(`select id, email from %s where id = $1`, UsersTable)
	var user entity.User
	err := r.db.GetContext(ctx, &user, q, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) GetProfileByUserId(ctx context.Context, userId int) (*entity.Profile, error) {
	q := fmt.Sprintf(`
		select email, first_name, last_name, phone, avatar, telegram, instagram 
		from %s where id = $1`,
		UsersTable,
	)
	var profile entity.Profile
	err := r.db.GetContext(ctx, &profile, q, userId)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *UsersRepo) CreateUser(ctx context.Context, user entity.User) error {
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

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}
