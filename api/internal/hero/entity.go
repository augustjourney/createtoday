package hero

import (
	"encoding/json"
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
