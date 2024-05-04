package entity

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Group struct {
	ID          int              `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description sql.NullString   `json:"description" db:"description"`
	Settings    *json.RawMessage `json:"settings" db:"settings"`
	ProjectID   int              `json:"project_id" db:"project_id"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

type UserGroup struct {
	ID       int        `json:"id" db:"id"`
	GroupID  int        `json:"group_id" db:"group_id"`
	UserID   int        `json:"user_id" db:"user_id"`
	Status   string     `json:"status" db:"status"`
	JoinedAt time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time `json:"left_at" db:"left_at"`
	RemoveAt *time.Time `json:"remove_at" db:"remove_at"`
}

type ProductGroup struct {
	ID               int              `json:"id" db:"id"`
	GroupID          int              `json:"group_id" db:"group_id"`
	ProductID        int              `json:"product_id" db:"product_id"`
	AllLessonsAccess bool             `json:"all_lessons_access" db:"all_lessons_access"`
	NoAccessContent  *json.RawMessage `json:"no_access_content" db:"no_access_content"`
}
