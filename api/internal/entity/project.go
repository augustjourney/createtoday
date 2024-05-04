package entity

import "time"

type Project struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Domain    string    `json:"domain" db:"domain"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
}

type NewProject struct {
	Name    string `json:"name" db:"name"`
	Domain  string `json:"domain" db:"domain"`
	OwnerID int    `json:"owner_id" db:"owner_id"`
}
