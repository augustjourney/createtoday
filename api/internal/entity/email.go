package entity

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
