package hero

import (
	"context"
)

// TODO: refactor to small interfaces
type Storage interface {
	// users
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user User) (int64, error)
	FindUserById(ctx context.Context, id int) (*User, error)
	UpdateUserInfo(ctx context.Context, dto UpdateUserInfoDTO) error

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
	CompleteLesson(ctx context.Context, lessonSlug string, userId int) error

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

	// offers
	GetOfferForRegistration(ctx context.Context, slug string) (*OfferForRegistration, error)
	GetOfferForProcessing(ctx context.Context, slug string) (*OfferForProcessing, error)
	GetOfferGroups(ctx context.Context, offerId int64) ([]int64, error)

	// payments
	GetPayMethods(ctx context.Context, projectId int64) ([]PayMethod, error)
	GetPayMethod(ctx context.Context, payMethodId int64, projectId int64) (*PayIntegration, error)

	// orders
	CreateOrder(ctx context.Context, order NewOrder) (int64, error)
	UpdateOrderPaymentId(ctx context.Context, orderId int64, paymentId string) error
	FindOrderById(ctx context.Context, orderId int64) (*OrderForProcessing, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, status string, orderError OrderError, cardInfo OrderCardInfo) error

	// enrollments
	AddUserToGroups(ctx context.Context, userId int64, groupIds []int64) error

	// quiz comments
	GetQuizComments(ctx context.Context, solvedQuizId int64) ([]QuizComment, error)
	CreateQuizComment(ctx context.Context, dto NewQuizComment) (int64, error)
	UpdateQuizComment(ctx context.Context, dto UpdateQuizComment) error
	DeleteQuizComment(ctx context.Context, quizCommentId int64, authorId int64) error
}
