package entity

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
