package hero

import (
	"database/sql"
	"encoding/json"
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
	Name      sql.NullString  `json:"name" db:"name"`
	Type      string          `json:"type" db:"type"`
	Slug      string          `json:"slug" db:"slug"`
	Mime      sql.NullString  `json:"mime" db:"mime"`
	Ext       sql.NullString  `json:"ext" db:"ext"`
	Size      *int            `json:"size" db:"size"`
	Width     *int            `json:"width" db:"width"`
	Height    *int            `json:"height" db:"height"`
	URL       sql.NullString  `json:"url" db:"url"`
	Sources   json.RawMessage `json:"sources" db:"sources"`
	Blurhash  json.RawMessage `json:"blurhash" db:"blurhash"`
	ParentID  *int            `json:"parent_id" db:"parent_id"`
	Storage   string          `json:"storage" db:"storage"`
	Duration  *int            `json:"duration" db:"duration"`
	Bucket    sql.NullString  `json:"bucket" db:"bucket"`
	Status    sql.NullString  `json:"status" db:"status"`
	Caption   sql.NullString  `json:"caption" db:"caption"`
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
	ID         int             `json:"id" db:"id"`
	UserID     int             `json:"user_id" db:"user_id"`
	QuizID     int             `json:"quiz_id" db:"quiz_id"`
	ProductID  int             `json:"product_id" db:"product_id"`
	LessonID   int             `json:"lesson_id" db:"lesson_id"`
	ProjectID  int             `json:"project_id" db:"project_id"`
	UserAnswer json.RawMessage `json:"user_answer" db:"user_answer"`
	Type       string          `json:"type" db:"type"`
	Starred    bool            `json:"starred" db:"starred"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

type QuizSolvedInfo struct {
	ID         int             `json:"id" db:"id"`
	UserAnswer json.RawMessage `json:"user_answer" db:"user_answer"`
	Type       string          `json:"type" db:"type"`
	Starred    bool            `json:"starred" db:"starred"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	Media      json.RawMessage `json:"media" db:"media"`
	Lesson     json.RawMessage `json:"lesson" db:"lesson"`
	Author     json.RawMessage `json:"author" db:"author"`
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
