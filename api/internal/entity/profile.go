package entity

type Profile struct {
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Phone     string `json:"phone" db:"phone"`
	Avatar    string `json:"avatar" db:"avatar"`
	Telegram  string `json:"telegram" db:"telegram"`
	Instagram string `json:"instagram" db:"instagram"`
}
