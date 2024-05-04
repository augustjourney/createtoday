package hero

import (
	"context"
	"createtodayapi/internal/common"
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
const LessonsTable = "public.lesson"

type PostgresRepo struct {
	db *sqlx.DB
}

func (r *PostgresRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	q := fmt.Sprintf(`select id, email, password from %s where email = $1`, UsersTable)
	var user User
	err := r.db.GetContext(ctx, &user, q, email)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) FindUserById(ctx context.Context, id int) (*User, error) {
	q := fmt.Sprintf(`select id, email from %s where id = $1`, UsersTable)
	var user User
	err := r.db.GetContext(ctx, &user, q, id)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) GetProfileByUserId(ctx context.Context, userId int) (*Profile, error) {
	q := fmt.Sprintf(`
		select email, first_name, last_name, phone, avatar, telegram, instagram 
		from %s where id = $1`,
		UsersTable,
	)
	var profile Profile
	err := r.db.GetContext(ctx, &profile, q, userId)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user User) error {
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

func (r *PostgresRepo) GetUserAccessibleProducts(ctx context.Context, userId int) ([]ProductCard, error) {
	q := fmt.Sprintf(`
		select distinct on (id) name, description, slug, cover, settings 
		from %s 
		where user_id = $1 and parent_id is null;
	`, UsersProductsView)
	var products []ProductCard
	err := r.db.SelectContext(ctx, &products, q, userId)
	if err != nil {
		return products, err
	}
	if products == nil {
		return []ProductCard{}, nil
	}
	return products, nil
}

func (r *PostgresRepo) GetUserAccessibleProduct(ctx context.Context, productSlug string, userId int) (*ProductInfo, error) {
	q := fmt.Sprintf(`
		select p.id, p.name, p.description, p.slug, p.cover, p.settings, p.layout 
		from %s as p
		where p.slug = $1 and p.user_id = $2
		limit 1
	`, UsersProductsView)
	var product ProductInfo
	err := r.db.GetContext(ctx, &product, q, productSlug, userId)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, common.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *PostgresRepo) GetProductLessons(ctx context.Context, productId int) ([]LessonCard, error) {
	q := fmt.Sprintf(`
		select l.name, l.slug, l.description
		from %s as l
		where l.product_id = $1 and l.is_published = true
		order by l.position asc
	`, LessonsTable)
	var lessons []LessonCard
	err := r.db.SelectContext(ctx, &lessons, q, productId)
	if err != nil {
		return make([]LessonCard, 0), err
	}
	if lessons == nil {
		return make([]LessonCard, 0), nil
	}
	return lessons, nil
}

func (r *PostgresRepo) GetUserAccessibleLesson(ctx context.Context, lessonSlug string, userId int) (*LessonInfo, error) {
	q := fmt.Sprintf(`
		select l.name, l.slug, l.description, l.content, l.can_complete,
		l.settings, l.is_public, json_build_object(
		    'name', p.name,
		    'slug', p.slug,
		    'cover', p.cover,
		    'settings', p.settings
		) as product
		from %s as l
		join %s as p on p.id = l.product_id and p.user_id = $2
		where l.slug = $1 and l.is_published = true
	`, LessonsTable, UsersProductsView)
	var lesson LessonInfo
	err := r.db.GetContext(ctx, &lesson, q, lessonSlug, userId)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrLessonNotFound
		}
		return nil, err
	}

	return &lesson, nil
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}
