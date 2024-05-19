package hero

import (
	"createtodayapi/internal/common"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *LoginBody) Validate() error {

	if b.Email == "" {
		return common.ErrEmptyEmail
	}

	if b.Password == "" {
		return common.ErrEmptyPassword
	}

	return nil
}

type GetMagicLinkBody struct {
	Email string `json:"email"`
}

func (b *GetMagicLinkBody) Validate() error {
	if b.Email == "" {
		return common.ErrEmptyEmail
	}
	return nil
}

type ValidateMagicLinkBody struct {
	Token string `json:"token"`
}

func (b *ValidateMagicLinkBody) Validate() error {
	if b.Token == "" {
		return common.ErrInvalidToken
	}

	return nil
}

type LoginResult struct {
	Token string `json:"token"`
}

type SignupBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
}

func (b *SignupBody) Validate() error {

	if b.Email == "" {
		return common.ErrEmptyEmail
	}

	return nil
}

type SignUpResult struct {
	AlreadyExists bool    `json:"alreadyExists"`
	Token         *string `json:"token"`
	Message       string  `json:"message,omitempty"`
}

type UpdateProfileBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Telegram  *string `json:"telegram"`
	Instagram *string `json:"instagram"`
	About     *string `json:"about"`
	Phone     *string `json:"phone"`
}

type UpdatePasswordBody struct {
	Password string `json:"password"`
}

func (b *UpdatePasswordBody) Validate() error {
	if b.Password == "" {
		return common.ErrNewPasswordIsEmpty
	}

	if len(b.Password) < 8 {
		return common.ErrNewPasswordIsShort
	}

	return nil
}

type SolveQuizBody struct {
	Answer string `json:"answer"`
	Type   string `json:"type"`
	Media  []FileUpload
	Slug   string
}

type FileUpload struct {
	FileName  string `json:"file_name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	FileURL   string `json:"file_url"`
	Mime      string `json:"mime"`
	MediaType string `json:"media_type"`
}

type FileUploadResult struct {
	MediaId int64  `json:"media_id"`
	FileURL string `json:"file_url"`
}

func (b *SolveQuizBody) Validate() error {
	if b.Type == "" {
		return common.ErrEmptyQuizType
	}

	if b.Type == "answer_quiz" && b.Answer == "" {
		return common.ErrEmptyQuizAnswer
	}

	if len(b.Media) == 0 && b.Type == "one_photo" {
		return common.ErrEmptyQuizPhoto
	}

	if len(b.Media) == 0 && b.Type == "one_video" {
		return common.ErrEmptyQuizVideo
	}

	if b.Type == "one_photo" && len(b.Media) > 1 {
		return common.ErrQuizTooManyPhotos
	}

	if b.Type == "one_video" && len(b.Media) > 1 {
		return common.ErrQuizTooManyVideos
	}

	return nil
}

type SolveQuizDTO struct {
	Answer   string
	Type     string
	UserID   int
	QuizSlug string
	Media    []FileUpload
}

type CreateUserDTO struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// TODO: сделать валидацию
// Получить оффер из бд
// Проверить обязательные поля — ask_for_telegram, ask_for_instagram, ask_for_phone
// По умолчанию обязательные поля — email, first_name, selected_pay_method
type ProcessOfferDTO struct {
	Slug              string `json:"slug"`
	UserID            int64  `json:"user_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Telegram          string `json:"telegram" db:"telegram"`
	Instagram         string `json:"instagram" db:"instagram"`
	Comment           string `json:"comment"`
	SelectedPayMethod int64  `json:"selected_pay_method" db:"selected_pay_method"`
}

type ProcessOfferResult struct {
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url"`
}

type CreatePaymentDTO struct {
	PayMethod        *PayIntegration
	UserID           int64
	Email            string
	Phone            string
	OrderDescription string
	OfferID          int64
	Price            uint64
}

type UpdateUserInfoDTO struct {
	UserID    int64  `json:"user_id" db:"user_id"`
	Phone     string `json:"phone" db:"phone"`
	Telegram  string `json:"telegram" db:"telegram"`
	Instagram string `json:"instagram" db:"instagram"`
}
