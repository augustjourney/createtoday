package hero

import (
	"context"
)

var emails []Email

var magicLinkLetter = Email{
	Subject:  "Cсылка для входа в CreateToday",
	Template: "default",
	From: EmailSender{
		Email: "hello@createtoday.ru",
		Name:  "CreateToday",
	},
	Body: `
		<p>
			Мы получили попытку входа в твой личный кабинет на {{ .Context.Domain }}.
			Чтобы войти — просто нажми на кнопку ниже.
		</p>
		<a href='{{ .Context.MagicLink }}' target='_blank' rel='noreferrer noopener' class='btn'>
			Войти
		</a>
		<p>Или используй ссылку:</p>
		<p>
			<a target='_blank' rel='noreferrer noopener' href='{{ .Context.MagicLink }}'>
				{{ .Context.MagicLink }}
			</a>
		</p>
		<p style='color: #475569;'>
			Если это не твоя попытка входа — просто проигнорируй это письмо.
		</p>
	`,
	IsActive: true,
	Type:     "magiclink",
	Context: map[string]interface{}{
		"Domain": "hero.createtoday.ru",
	},
}

var welcomeLetter = Email{
	Subject:  "Добро пожаловать в CreateToday",
	Template: "default",
	From: EmailSender{
		Email: "hello@createtoday.ru",
		Name:  "CreateToday",
	},
	IsActive: true,
	Type:     "welcome",
	Context: map[string]interface{}{
		"Domain": "hero.createtoday.ru",
	},
	Body: `
		<p>Привет! 👋</p>
		<p>Для тебя создана учетная запись на сайте {{ .Context.Domain }}.</p>
		<p>Вот данные для входа:</p>
		<p style='margin-bottom: 0'>
			Логин: <span style='color: #0284c7;'>{{ .Context.Email }}</span>
		</p>
		<p>
			Пароль: <span style='color: #0284c7;'>{{ .Context.Password }}</span>
		</p>
		<a href='{{ .Context.LoginFullURL }}' target='_blank' rel='noreferrer noopener' class='btn'>
			Войти
		</a>
		<p>
			Чтобы войти в свой личный кабинет, нажми на кнопку выше. Или используйте ссылку
			<a target='_blank' rel='noreferrer noopener' href='{{ .Context.LoginFullURL }}'>
				{{ .Context.LoginURL }}
			</a>
		</p>
		<p>
			Пароль сможешь поменять в своем личном кабинете. А если есть вопросы —
			вот наша почта: <span style='color: #0284c7;'>{{ .Context.MailFrom }}</span>
		</p>
	`,
}

type MemoryRepo struct{}

func (r *MemoryRepo) FindByType(ctx context.Context, emailType string) (*Email, error) {
	var email Email
	for _, e := range emails {
		if e.Type == emailType {
			email = e
		}
	}
	return &email, nil
}

func NewMemoryRepo() *MemoryRepo {
	emails = append(emails, magicLinkLetter, welcomeLetter)
	return &MemoryRepo{}
}
