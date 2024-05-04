package entity

import (
	"encoding/json"
	"time"
)

type ProductSettings struct {
	ContentTitle     string `json:"content_title"`
	SubproductsTitle string `json:"subproducts_title"`
	Cover            struct {
		ShowGradient    bool `json:"show_gradient"`
		ShowDescription bool `json:"show_description"`
	} `json:"cover"`
}

type ProductCover struct {
	MediaId int    `json:"media_id"`
	Url     string `json:"url"`
	Sources []struct {
		Mime   string `json:"mime"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Image  string `json:"image"`
		Url    string `json:"url"`
	} `json:"sources"`
}

type Product struct {
	ID                       int             `json:"id" db:"id"`
	Name                     string          `json:"name" db:"name"`
	Slug                     string          `json:"slug" db:"slug"`
	Description              *string         `json:"description" db:"description"`
	Layout                   string          `json:"layout" db:"layout"`
	Position                 int             `json:"position" db:"position"`
	Access                   string          `json:"access" db:"access"`
	IsPublished              bool            `json:"is_published" db:"is_published"`
	Cover                    ProductCover    `json:"cover" db:"cover"`
	ParentID                 *int            `json:"parent_id" db:"parent_id"`
	ProjectID                int             `json:"project_id" db:"project_id"`
	Settings                 ProductSettings `json:"settings" db:"settings"`
	ShowLessonsWithoutAccess bool            `json:"show_lessons_without_access" db:"show_lessons_without_access"`
	CreatedAt                time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at" db:"updated_at"`
	CreatedBy                int             `json:"created_by" db:"created_by"`
}

type NewProduct struct {
	Name      string `json:"name" db:"name"`
	Slug      string `json:"slug" db:"slug"`
	ProjectID int    `json:"project_id" db:"project_id"`
	CreatedBy int    `json:"created_by" db:"created_by"`
	Position  int    `json:"position" db:"position"`
}

type UserProduct struct {
	ID                       int             `json:"id" db:"id"`
	Name                     string          `json:"name" db:"name"`
	Slug                     string          `json:"slug" db:"slug"`
	Description              string          `json:"description" db:"description"`
	Layout                   string          `json:"layout" db:"layout"`
	Position                 int             `json:"position" db:"position"`
	IsPublished              bool            `json:"is_published" db:"is_published"`
	Cover                    ProductCover    `json:"cover" db:"cover"`
	ParentID                 *int            `json:"parent_id" db:"parent_id"`
	ProjectID                int             `json:"project_id" db:"project_id"`
	Settings                 ProductSettings `json:"settings" db:"settings"`
	ShowLessonsWithoutAccess bool            `json:"show_lessons_without_access" db:"show_lessons_without_access"`
}

type UserProductCard struct {
	Name        string           `json:"name" db:"name"`
	Slug        string           `json:"slug" db:"slug"`
	Description *string          `json:"description" db:"description"`
	Cover       *json.RawMessage `json:"cover" db:"cover"`
	Settings    *json.RawMessage `json:"settings" db:"settings"`
}
