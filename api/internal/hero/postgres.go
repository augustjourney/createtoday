package hero

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/logger"
	"database/sql"
	"errors"
	"fmt"

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
const QuizzesTable = "public.quiz"
const MediaTable = "public.media"
const RelatedMediaTable = "public.related_media"
const SolvedQuizzesTable = "public.quiz_solved"
const SolvedQuizzesView = "public._solvedquizzes"

type PostgresRepo struct {
	db *sqlx.DB
}

func (r *PostgresRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	q := fmt.Sprintf(`select id, email, password from %s where email = $1`, UsersTable)
	var user User
	err := r.db.GetContext(ctx, &user, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) GetProfileByUserId(ctx context.Context, userId int) (*Profile, error) {
	q := fmt.Sprintf(`
		select email, first_name, last_name, phone, avatar, telegram, instagram, about 
		from %s where id = $1`,
		UsersTable,
	)
	var profile Profile
	err := r.db.GetContext(ctx, &profile, q, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *PostgresRepo) UpdateProfile(ctx context.Context, userId int, profile UpdateProfileBody) error {
	q := fmt.Sprintf(`
		update %s 
		set first_name = $2, last_name = $3, phone = $4, telegram = $5, instagram = $6, about = $7
		where id = $1`,
		UsersTable,
	)
	_, err := r.db.ExecContext(ctx, q, userId, profile.FirstName, profile.LastName, profile.Phone, profile.Telegram, profile.Instagram, profile.About)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) UpdateAvatar(ctx context.Context, userId int, avatar string) error {
	q := fmt.Sprintf(`
		update %s set avatar = $2
		where id = $1`,
		UsersTable,
	)
	_, err := r.db.ExecContext(ctx, q, userId, avatar)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) UpdatePassword(ctx context.Context, userId int, password string) error {
	q := fmt.Sprintf(`
		update %s set password = $2
		where id = $1`,
		UsersTable,
	)
	_, err := r.db.ExecContext(ctx, q, userId, password)
	if err != nil {
		return err
	}
	return nil
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
		if errors.Is(err, sql.ErrNoRows) {
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
		l.settings, l.is_public, 
		json_build_object(
		    'name', p.name,
		    'slug', p.slug,
		    'cover', p.cover,
		    'settings', p.settings
		) as product, 
		coalesce(lq.quizzes, '{}'::json) as quizzes,
		coalesce(lm.media, '{}'::json) as media
		
		from %s as l
		join %s as p on p.id = l.product_id and p.user_id = $2
		
		-- quizzes
		left join lateral (
		    select 
		    	json_object_agg(
		    		q.id,
		    		json_build_object(
		    			'slug', q.slug,
		    			'type', q.type,
		    			'content', q.content,
		    			'settings', q.settings,
		    			'show_others_answers', q.show_others_answers
		    		)
		    	) as quizzes
		    from %s as q
		    where q.lesson_id = l.id
		) as lq on true
		
		-- media
		left join lateral (
		    select 
		    	json_object_agg(
		    		rm.media_id,
		    		json_build_object(
		    			'url', m.url,
		    			'sources', m.sources,
		    			'type', m.type
		    		)
		    	) as media
		    from %s as rm
		    join %s as m on m.id = rm.media_id
		    where rm.related_type = 'lesson' and rm.related_id = l.id
		) as lm on true
		             
		where l.slug = $1 and l.is_published = true
	`, LessonsTable, UsersProductsView, QuizzesTable, RelatedMediaTable, MediaTable)

	var lesson LessonInfo

	err := r.db.GetContext(ctx, &lesson, q, lessonSlug, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrLessonNotFound
		}
		return nil, err
	}

	return &lesson, nil
}

func (r *PostgresRepo) GetSolvedQuizzesForProduct(ctx context.Context, productSlug string, userId int) ([]QuizSolvedInfo, error) {
	var solvedQuizzes []QuizSolvedInfo

	q := fmt.Sprintf(`
		select q.id, q.user_answer, q.type, q.created_at, q.starred, q.media, q.lesson, q.author
		from %s as q
		where q.product_id = (select id from %s where slug = $1 and user_id = $2)		
		order by q.created_at desc;
	`, SolvedQuizzesView, UsersProductsView)

	err := r.db.SelectContext(ctx, &solvedQuizzes, q, productSlug, userId)
	if err != nil {
		return make([]QuizSolvedInfo, 0), err
	}

	if solvedQuizzes == nil {
		return make([]QuizSolvedInfo, 0), nil
	}

	return solvedQuizzes, nil
}

func (r *PostgresRepo) GetSolvedQuizzesForUser(ctx context.Context, productSlug string, userId int) ([]QuizSolvedInfo, error) {
	var solvedQuizzes []QuizSolvedInfo

	q := fmt.Sprintf(`
		select q.id, q.user_answer, q.type, q.created_at, q.starred, q.media, q.lesson, q.author
		from %s as q
		where q.product_id = (select id from %s where slug = $1 and user_id = $2) and user_id = $2	
		order by q.created_at desc;
	`, SolvedQuizzesView, UsersProductsView)

	err := r.db.SelectContext(ctx, &solvedQuizzes, q, productSlug, userId)
	if err != nil {
		return make([]QuizSolvedInfo, 0), err
	}

	if solvedQuizzes == nil {
		return make([]QuizSolvedInfo, 0), nil
	}

	return solvedQuizzes, nil
}

func (r *PostgresRepo) GetSolvedQuizzesForQuiz(ctx context.Context, quizSlug string) ([]QuizSolvedInfo, error) {
	var solvedQuizzes []QuizSolvedInfo

	q := fmt.Sprintf(`
		select q.id, q.user_answer, q.type, q.created_at, q.starred, q.media, q.lesson, q.author
		from %s as q
		where q.quiz_id = (select id from %s where slug = $1)
		
		order by q.created_at desc;
	`, SolvedQuizzesView, QuizzesTable)

	err := r.db.SelectContext(ctx, &solvedQuizzes, q, quizSlug)
	if err != nil {
		return make([]QuizSolvedInfo, 0), err
	}

	if solvedQuizzes == nil {
		return make([]QuizSolvedInfo, 0), nil
	}

	return solvedQuizzes, nil
}

func (r *PostgresRepo) SolveQuiz(ctx context.Context, quizSlug string, userId int, answer []byte) (int64, error) {
	q := fmt.Sprintf(`
		insert into %s (user_id, quiz_id, product_id, lesson_id, project_id, type, user_answer)
		select :user_id, id, product_id, lesson_id, project_id, type, :answer
		from %s 
		where slug = :slug
		returning id
	`, SolvedQuizzesTable, QuizzesTable)

	type solveQuizArgs struct {
		Slug   string `json:"slug" db:"slug"`
		UserId int    `json:"user_id" db:"user_id"`
		Answer []byte `json:"answer" db:"answer"`
	}

	quizArgs := solveQuizArgs{
		Slug:   quizSlug,
		UserId: userId,
		Answer: answer,
	}

	query, args, err := r.db.BindNamed(q, quizArgs)

	if err != nil {
		return 0, err
	}

	var solvedQuizId int64

	err = r.db.GetContext(ctx, &solvedQuizId, query, args...)
	if err != nil {
		return 0, err
	}

	return solvedQuizId, nil
}

func (r *PostgresRepo) SaveMedia(ctx context.Context, media Media) (int64, error) {

	q := fmt.Sprintf(`
		insert into %s
		(name, slug, bucket, mime, ext, storage, original, width, height, size, duration, type, url, status)
		values (:name, :slug, :bucket, :mime, :ext, :storage, :original, :width, :height, :size, :duration, :type, :url, :status)
		returning id;
	`, MediaTable)

	query, args, err := r.db.BindNamed(q, media)

	if err != nil {
		logger.Log.Error(err.Error(), "where", "SaveMedia.bindNamed()")
		return 0, err
	}

	var mediaId int64

	err = r.db.GetContext(ctx, &mediaId, query, args...)
	if err != nil {
		logger.Log.Error(err.Error(), "where", "SaveMedia.GetContext()")
		return 0, err
	}

	return mediaId, nil
}

func (r *PostgresRepo) ConnectManyMedia(ctx context.Context, mediaIds []int64, relatedType string, relatedId int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, mediaId := range mediaIds {
		q := fmt.Sprintf(`
			insert into %s
			(media_id, related_type, related_id)
			values ($1, $2, $3)
		`, RelatedMediaTable)
		_, err := tx.ExecContext(ctx, q, mediaId, relatedType, relatedId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepo) ConnectMedia(ctx context.Context, mediaId int64, relatedType string, relatedId int64) error {
	q := fmt.Sprintf(`
			insert into %s
			(media_id, related_type, related_id)
			values ($1, $2, $3)
		`, RelatedMediaTable)

	_, err := r.db.ExecContext(ctx, q, mediaId, relatedType, relatedId)

	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.Postgres.ConnectMedia")
		return err
	}

	return nil
}

func (r *PostgresRepo) DeleteMedia(ctx context.Context, mediaId int64) error {
	q := fmt.Sprintf(`
		update %s set status = 'to_delete' where id = $1; 
		delete from %s where media_id = $1
	`, MediaTable, RelatedMediaTable)

	_, err := r.db.ExecContext(ctx, q, mediaId)
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.DeleteMedia")
		return err
	}
	return err
}

func (r *PostgresRepo) UpdateMediaStatus(ctx context.Context, mediaId int64, status string) error {
	q := fmt.Sprintf(`update %s set status = $2 where id = $1`, MediaTable)
	_, err := r.db.ExecContext(ctx, q, mediaId, status)
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.UpdateMediaStatus")
		return err
	}
	return nil
}

func (r *PostgresRepo) FindSolvedQuiz(ctx context.Context, userId int, quizSlug string) (*QuizSolved, error) {
	q := fmt.Sprintf(`
		select * from %s 
		where user_id = $1
		and quiz_id = (select id from %s where slug = $2)
	`, SolvedQuizzesTable, QuizzesTable)

	var quizSolved QuizSolved

	err := r.db.GetContext(ctx, &quizSolved, q, userId, quizSlug)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, common.ErrSolvedQuizNotFound
	}

	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.FindSolvedQuiz")
		return nil, err
	}

	return &quizSolved, nil
}

func (r *PostgresRepo) GetQuizBySlug(ctx context.Context, quizSlug string) (*Quiz, error) {
	q := fmt.Sprintf(`select * from %s where slug = $1`, QuizzesTable)

	var quiz Quiz
	err := r.db.GetContext(ctx, &quiz, q, quizSlug)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, common.ErrQuizNotFound
	}

	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.GetQuizBySlug")
		return nil, err
	}

	return &quiz, nil
}

func (r *PostgresRepo) DeleteSolvedQuiz(ctx context.Context, solvedQuizId int64, userId int) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q1 := fmt.Sprintf(`
		update %s set status = $2 
		where id = (select media_id from %s where related_type = 'solved_quiz' and related_id = $1)
	`, MediaTable, RelatedMediaTable)

	_, err = tx.ExecContext(ctx, q1, solvedQuizId, "to_delete")
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.DeleteSolvedQuiz.Q1")
		_ = tx.Rollback()
		return err
	}

	q2 := fmt.Sprintf(`
		delete from %s where related_type = 'solved_quiz' and related_id = $1
	`, RelatedMediaTable)

	_, err = tx.ExecContext(ctx, q2, solvedQuizId)
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.DeleteSolvedQuiz.Q2")
		_ = tx.Rollback()
		return err
	}

	q3 := fmt.Sprintf(`
		delete from %s where id = $1 and user_id = $2
	`, SolvedQuizzesTable)

	_, err = tx.ExecContext(ctx, q3, solvedQuizId, userId)
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.DeleteSolvedQuiz.Q3")
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error(), "where", "hero.postgres.DeleteSolvedQuiz.Commit")
		return err
	}

	return nil
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}
