package entity

import (
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
	IsDeleted    bool             `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" db:"updated_at"`
	CreatedBy    int              `json:"created_by" db:"created_by"`
	ProjectID    int              `json:"project_id" db:"project_id"`
	ProductID    int              `json:"product_id" db:"product_id"`
	PublishAt    *time.Time       `json:"publish_at" db:"publish_at"`
}
