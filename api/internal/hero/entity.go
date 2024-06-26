package hero

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Lesson struct {
	ID           int              `json:"id" db:"id"`
	Name         string           `json:"name" db:"name"`
	Slug         string           `json:"slug" db:"slug"`
	Description  *string          `json:"description" db:"description"`
	Content      json.RawMessage  `json:"content" db:"content"`
	Position     int              `json:"position" db:"position"`
	IsPublished  bool             `json:"is_published" db:"is_published"`
	IsStopLesson bool             `json:"is_stop_lesson" db:"is_stop_lesson"`
	CanComplete  bool             `json:"can_complete" db:"can_complete"`
	Settings     *json.RawMessage `json:"settings" db:"settings"`
	IsPublic     bool             `json:"is_public" db:"is_public"`
}

type Profile struct {
	Email     string  `json:"email" db:"email"`
	FirstName *string `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name" db:"last_name"`
	Phone     *string `json:"phone" db:"phone"`
	Avatar    *string `json:"avatar" db:"avatar"`
	Telegram  *string `json:"telegram" db:"telegram"`
	Instagram *string `json:"instagram" db:"instagram"`
	About     *string `json:"about" db:"about"`
}

type User struct {
	ID        int    `json:"id" db:"id"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Phone     string `json:"phone" db:"phone"`
	Avatar    string `json:"avatar" db:"avatar"`
	Telegram  string `json:"telegram" db:"telegram"`
	Instagram string `json:"instagram" db:"instagram"`
	LastSeen  string `json:"last_seen" db:"last_seen"`
}

type ProductCard struct {
	Name        string           `json:"name" db:"name"`
	Slug        string           `json:"slug" db:"slug"`
	Description *string          `json:"description" db:"description"`
	Cover       *json.RawMessage `json:"cover" db:"cover"`
	Settings    *json.RawMessage `json:"settings" db:"settings"`
}

type EmailSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Email struct {
	Subject   string
	Body      string
	Type      string
	From      EmailSender
	Template  string
	Context   map[string]interface{}
	ProjectID int
	IsActive  bool
}

type Product struct {
	ID                       int              `json:"id" db:"id"`
	Name                     string           `json:"name" db:"name"`
	Slug                     string           `json:"slug" db:"slug"`
	Description              string           `json:"description" db:"description"`
	Layout                   string           `json:"layout" db:"layout"`
	Position                 int              `json:"position" db:"position"`
	IsPublished              bool             `json:"is_published" db:"is_published"`
	Cover                    *json.RawMessage `json:"cover" db:"cover"`
	ParentID                 *int             `json:"parent_id" db:"parent_id"`
	ProjectID                int              `json:"project_id" db:"project_id"`
	Settings                 *json.RawMessage `json:"settings" db:"settings"`
	ShowLessonsWithoutAccess bool             `json:"show_lessons_without_access" db:"show_lessons_without_access"`
}

type ProductInfo struct {
	ID                       int              `json:"-" db:"id"`
	Name                     string           `json:"name" db:"name"`
	Slug                     string           `json:"slug" db:"slug"`
	Description              *string          `json:"description" db:"description"`
	Layout                   string           `json:"layout" db:"layout"`
	Cover                    *json.RawMessage `json:"cover" db:"cover"`
	ProjectID                int              `json:"-" db:"project_id"`
	Settings                 *json.RawMessage `json:"settings" db:"settings"`
	ShowLessonsWithoutAccess bool             `json:"show_lessons_without_access" db:"show_lessons_without_access"`
	Lessons                  []LessonCard     `json:"lessons"`
}

type LessonCard struct {
	Name        string  `json:"name" db:"name"`
	Slug        string  `json:"slug" db:"slug"`
	Description *string `json:"description" db:"description"`
}

type LessonInfo struct {
	ID          int              `json:"-" db:"id"`
	Name        string           `json:"name" db:"name"`
	Slug        string           `json:"slug" db:"slug"`
	Description *string          `json:"description" db:"description"`
	Content     *json.RawMessage `json:"content" db:"content"`
	CanComplete bool             `json:"can_complete" db:"can_complete"`
	Settings    *json.RawMessage `json:"settings" db:"settings"`
	IsPublic    bool             `json:"-" db:"is_public"`
	Product     json.RawMessage  `json:"product" db:"product"`
	Quizzes     json.RawMessage  `json:"quizzes" db:"quizzes"`
	Media       json.RawMessage  `json:"media" db:"media"`
	NextLesson  *string          `json:"next_lesson" db:"next_lesson"`
}

type LessonContent struct {
	Elements []LessonElement `json:"elements"`
}

type LessonElement struct {
	ID   string      `json:"id" db:"id"`
	Type string      `json:"type" db:"type"`
	Body interface{} `json:"body" db:"body"`
}

type LessonElementGallery struct {
	Media []struct {
		MediaId int `json:"media_id"`
		MediaInfo
	} `json:"media"`
	Settings struct {
		View string `json:"view"`
	} `json:"settings"`
}

type LessonElementQuiz struct {
	QuizId int `json:"quiz_id"`
	QuizInfo
}

type LessonElementGif struct {
	Media struct {
		MediaId int `json:"media_id"`
		MediaInfo
	} `json:"media"`
	Settings *json.RawMessage `json:"settings"`
}

type LessonElementAudio struct {
	Media struct {
		MediaId int `json:"media_id"`
		MediaInfo
	} `json:"media"`
	Settings *json.RawMessage `json:"settings"`
}

type Quiz struct {
	ID                int             `json:"id" db:"id"`
	Name              sql.NullString  `json:"name" db:"name"`
	Slug              string          `json:"slug" db:"slug"`
	Content           json.RawMessage `json:"content" db:"content"`
	Type              string          `json:"type" db:"type"`
	Settings          json.RawMessage `json:"settings" db:"settings"`
	ShowOthersAnswers bool            `json:"show_others_answers" db:"show_others_answers"`
	LessonID          int             `json:"lesson_id" db:"lesson_id"`
	CreatedBy         *int            `json:"created_by" db:"created_by"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	ProjectID         int             `json:"project_id" db:"project_id"`
	ProductID         int             `json:"product_id" db:"product_id"`
}

type QuizContent struct {
	Caption string `json:"caption" db:"caption"`
}

type QuizSettings struct {
	StopBlock bool `json:"stop_block" db:"stop_block"`
}

// Как используется в уроке
type QuizInfo struct {
	Slug              string       `json:"slug" db:"slug"`
	Content           QuizContent  `json:"content" db:"content"`
	Type              string       `json:"type" db:"type"`
	Settings          QuizSettings `json:"settings" db:"settings"`
	ShowOthersAnswers bool         `json:"show_others_answers" db:"show_others_answers"`
}

type Media struct {
	ID        int             `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Type      string          `json:"type" db:"type"`
	Slug      string          `json:"slug" db:"slug"`
	Mime      string          `json:"mime" db:"mime"`
	Ext       string          `json:"ext" db:"ext"`
	Size      *int64          `json:"size" db:"size"`
	Width     *int            `json:"width" db:"width"`
	Height    *int            `json:"height" db:"height"`
	URL       string          `json:"url" db:"url"`
	Sources   json.RawMessage `json:"sources" db:"sources"`
	Blurhash  json.RawMessage `json:"blurhash" db:"blurhash"`
	ParentID  *int            `json:"parent_id" db:"parent_id"`
	Storage   string          `json:"storage" db:"storage"`
	Duration  *int            `json:"duration" db:"duration"`
	Bucket    string          `json:"bucket" db:"bucket"`
	Status    string          `json:"status" db:"status"`
	Caption   *string         `json:"caption" db:"caption"`
	Original  bool            `json:"original" db:"original"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type MediaSource struct {
	Mime   string `json:"mime" db:"mime"`
	Width  int    `json:"width" db:"width"`
	Height int    `json:"height" db:"height"`
	Image  string `json:"image" db:"image"` // image name
	Url    string `json:"url" db:"url"`
}

// Как используется в уроках и выполненных квизах
type MediaInfo struct {
	Type    string        `json:"type" db:"type"`
	Url     string        `json:"url" db:"url"`
	Sources []MediaSource `json:"sources" db:"sources"`
}

type QuizSolved struct {
	ID         int64           `json:"id" db:"id"`
	UserID     int             `json:"user_id" db:"user_id"`
	QuizID     int             `json:"quiz_id" db:"quiz_id"`
	ProductID  int             `json:"product_id" db:"product_id"`
	LessonID   int             `json:"lesson_id" db:"lesson_id"`
	ProjectID  int             `json:"project_id" db:"project_id"`
	UserAnswer json.RawMessage `json:"user_answer" db:"user_answer"`
	Type       string          `json:"type" db:"type"`
	Starred    bool            `json:"starred" db:"starred"`
	IsDeleted  bool            `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

type QuizSolvedAnswer struct {
	Answer string `json:"answer" db:"answer"`
}

type QuizSolvedInfo struct {
	ID         int              `json:"id" db:"id"`
	UserAnswer json.RawMessage  `json:"user_answer" db:"user_answer"`
	Type       string           `json:"type" db:"type"`
	Starred    bool             `json:"starred" db:"starred"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
	Media      *json.RawMessage `json:"media" db:"media"`
	Lesson     json.RawMessage  `json:"lesson" db:"lesson"`
	Author     json.RawMessage  `json:"author" db:"author"`
}

type QuizSolvedAuthor struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Avatar    string `json:"avatar" db:"avatar"`
}

type QuizSolvedLesson struct {
	Name    string `json:"name" db:"name"`
	Slug    string `json:"slug" db:"slug"`
	Product struct {
		Name string `json:"name" db:"name"`
		Slug string `json:"slug" db:"slug"`
	} `json:"product" db:"product"`
}

type Offer struct {
	ID                     int64            `db:"id"`
	Name                   string           `db:"name"`
	Description            *string          `db:"description"`
	Slug                   string           `db:"slug"`
	Price                  uint64           `db:"price"`
	Currency               string           `db:"currency"`
	IsFree                 bool             `db:"is_free"`
	SendOrderCreated       bool             `db:"send_order_created"`
	SendOrderCompleted     bool             `db:"send_order_completed"`
	SendRegistrationEmail  bool             `db:"send_registration_email"`
	RegistrationEmail      *string          `db:"registration_email"`
	AddToNewsletter        bool             `db:"add_to_newsletter"`
	Type                   string           `db:"type"`
	Settings               *json.RawMessage `db:"settings"`
	ProjectID              int              `db:"project_id"`
	CreatedAt              time.Time        `db:"created_at"`
	UpdatedAt              time.Time        `db:"updated_at"`
	RegistrationEmailTheme *string          `db:"registration_email_theme"`
	SuccessMessage         *string          `db:"success_message"`
	RedirectURL            *string          `db:"redirect_url"`
	AskForPhone            bool             `db:"ask_for_phone"`
	AskForComment          bool             `db:"ask_for_comment"`
	OfertaURL              *string          `db:"oferta_url"`
	AgreementURL           *string          `db:"agreement_url"`
	PrivacyURL             *string          `db:"privacy_url"`
	SendWelcomeEmail       bool             `db:"send_welcome_email"`
	AskForTelegram         bool             `db:"ask_for_telegram"`
	AskForInstagram        bool             `db:"ask_for_instagram"`
	CanUsePromocode        bool             `db:"can_use_promocode"`
	IsDonate               bool             `db:"is_donate"`
	MinDonatePrice         int              `db:"min_donate_price"`
	SendToSalebot          bool             `db:"send_to_salebot"`
	SalebotCallbackText    *string          `db:"salebot_callback_text"`
}

type OfferForRegistration struct {
	Name            string           `db:"name" json:"name"`
	Description     *string          `db:"description" json:"description"`
	Slug            string           `db:"slug" json:"slug"`
	Price           uint64           `db:"price" json:"price"`
	Currency        string           `db:"currency" json:"currency"`
	IsFree          bool             `db:"is_free" json:"is_free"`
	Settings        *json.RawMessage `db:"settings" json:"settings"`
	AskForPhone     bool             `db:"ask_for_phone" json:"ask_for_phone"`
	AskForComment   bool             `db:"ask_for_comment" json:"ask_for_comment"`
	OfertaURL       *string          `db:"oferta_url" json:"oferta_url"`
	AgreementURL    *string          `db:"agreement_url" json:"agreement_url"`
	PrivacyURL      *string          `db:"privacy_url" json:"privacy_url"`
	AskForTelegram  bool             `db:"ask_for_telegram" json:"ask_for_telegram"`
	AskForInstagram bool             `db:"ask_for_instagram" json:"ask_for_instagram"`
	CanUsePromocode bool             `db:"can_use_promocode" json:"can_use_promocode"`
	IsDonate        bool             `db:"is_donate" json:"is_donate"`
	MinDonatePrice  int              `db:"min_donate_price" json:"min_donate_price"`
	PayMethods      json.RawMessage  `db:"pay_methods" json:"pay_methods"`
	CanProcess      bool             `db:"can_process" json:"can_process"`
}

type OfferForProcessing struct {
	ID                     int64            `db:"id"`
	Name                   string           `db:"name"`
	Slug                   string           `db:"slug"`
	Price                  uint64           `db:"price"`
	Currency               string           `db:"currency"`
	IsFree                 bool             `db:"is_free"`
	SendOrderCreated       bool             `db:"send_order_created"`
	SendOrderCompleted     bool             `db:"send_order_completed"`
	SendRegistrationEmail  bool             `db:"send_registration_email"`
	RegistrationEmail      *string          `db:"registration_email"`
	Settings               *json.RawMessage `db:"settings"`
	ProjectID              int64            `db:"project_id"`
	RegistrationEmailTheme *string          `db:"registration_email_theme"`
	SuccessMessage         *string          `db:"success_message"`
	RedirectURL            *string          `db:"redirect_url"`
	SendWelcomeEmail       bool             `db:"send_welcome_email"`
	CanUsePromocode        bool             `db:"can_use_promocode"`
	IsDonate               bool             `db:"is_donate"`
	MinDonatePrice         int              `db:"min_donate_price"`
	SendToSalebot          bool             `db:"send_to_salebot"`
	SalebotCallbackText    *string          `db:"salebot_callback_text"`
	PayMethod              *PayIntegration  `db:"pay_method"`
}

type ReceiptSettings struct {
	Taxation string `json:"taxation"`
}

type PayIntegration struct {
	ID              int64            `json:"id" db:"id"`
	Name            string           `json:"name" db:"name"`
	Type            string           `json:"type" db:"type"`
	Login           string           `json:"login" db:"login"`
	Password        string           `json:"password" db:"password"`
	IsActive        bool             `json:"is_active" db:"is_active"`
	SendReceipt     bool             `json:"send_receipt" db:"send_receipt"`
	ReceiptSettings *ReceiptSettings `json:"receipt_settings" db:"receipt_settings"`
	ProjectID       int64            `json:"project_id" db:"project_id"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
}

type PayMethod struct {
	Name string `json:"name" db:"name"`
	Type string `json:"type" db:"type"`
}

type Order struct {
	ID              int64            `db:"id"`
	Description     *string          `db:"description"`
	Comment         *string          `db:"comment"`
	Price           uint64           `db:"price"`
	Currency        string           `db:"currency"`
	Status          string           `db:"status"`
	Error           *json.RawMessage `db:"error"`
	CardInfo        *json.RawMessage `db:"card_info"`
	PaymentID       string           `db:"payment_id"`
	IntegrationID   int64            `db:"integration_id"`
	OfferID         int64            `db:"offer_id"`
	ProjectID       int64            `db:"project_id"`
	UserID          int64            `db:"user_id"`
	Note            *string          `db:"note"`
	IsGift          bool             `db:"is_gift"`
	GiftedTo        *json.RawMessage `db:"gifted_to"`
	SalebotClientID *string          `db:"salebot_client_id"`
	CreatedAt       time.Time        `db:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at"`
}

type OrderForProcessing struct {
	ID        int64  `db:"id"`
	Price     uint64 `db:"price"`
	OfferID   int64  `db:"offer_id"`
	OfferSlug string `db:"offer_slug"`
	Status    string `db:"status"`
	UserID    int64  `db:"user_id"`
	UserEmail string `db:"user_email"`
	PaymentID string `db:"payment_id"`
}

type NewOrder struct {
	IntegrationID int64 `db:"integration_id"`
	OfferID       int64 `db:"offer_id"`
	UserID        int64 `db:"user_id"`
}

type OrderError struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

type OrderCardInfo struct {
	ExpirationDate string `json:"expiration_date"`
	Pan            string `json:"pan"`
}

type QuizComment struct {
	ID              int64             `db:"id" json:"id"`
	AuthorID        int64             `db:"-" json:"-"`
	QuizSolvedID    int64             `db:"-" json:"-"`
	UUID            string            `db:"uuid" json:"uuid"`
	Text            string            `db:"text" json:"text"`
	IsRead          bool              `db:"is_read" json:"is_read"`
	IsEdited        bool              `db:"is_edited" json:"is_edited"`
	IsFromModerator bool              `db:"is_from_moderator" json:"is_from_moderator"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
	Author          QuizCommentAuthor `db:"author" json:"author"`
}

type QuizCommentAuthor struct {
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Avatar    string `db:"avatar" json:"avatar"`
}

func (q *QuizCommentAuthor) Scan(v interface{}) error {
	switch vv := v.(type) {
	case []byte:
		return json.Unmarshal(vv, q)
	case string:
		return json.Unmarshal([]byte(vv), q)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type NewQuizComment struct {
	AuthorID     int64  `db:"author_id" json:"author_id"`
	QuizSolvedID int64  `db:"quiz_solved_id" json:"quiz_solved_id"`
	UUID         string `db:"uuid" json:"uuid"`
	Text         string `db:"text" json:"text"`
}
