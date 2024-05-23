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

var orderCreated = Email{
	Subject:  "Заказ принят",
	Template: "default",
	From: EmailSender{
		Email: "hello@createtoday.ru",
		Name:  "CreateToday",
	},
	IsActive: true,
	Type:     "order-created",
	Context: map[string]interface{}{
		"Domain":    "hero.createtoday.ru",
		"RespondTo": "hello@createtoday.ru",
	},
	Body: `
		<h3>Заказ принят! 🚀 </h3>

		<p>Привет! Только что поступил твой заказ на <strong>{{ .Context.Ordered }}</strong> на сумму <strong>{{ .Context.Amount }}₽</strong>.</p>

		<p>Скоро у тебя будет кое-что очень крутое 🤩. Для подтверждения заказа —
			оплати его по ссылке:</p>

		<a class='btn' target='_blank' rel='noreferrer noopener' href='{{ .Context.PaymentURL }}'>
			Перейти к оплате
		</a>

		<p>
			Если появятся вопросы, вот наша почта: {{ .Context.RespondTo }}.
		</p>

		<p>Успехов, <br />команда create.today</p>
	`,
}

var orderCompleted = Email{
	Subject:  "Заказ оплачен! 🥳",
	Template: "default",
	From: EmailSender{
		Email: "hello@createtoday.ru",
		Name:  "CreateToday",
	},
	IsActive: true,
	Type:     "order-completed",
	Context: map[string]interface{}{
		"Domain":    "hero.createtoday.ru",
		"RespondTo": "hello@createtoday.ru",
	},
	Body: `
		<h3>Заказ оплачен! 🥳</h3>
		<p>Спасибо за доверие!</p>
		<p style='margin-bottom: 10px !important;'>Твой заказ:
			<strong>{{ .Context.Ordered }}</strong>
		</p>
		<p>Сумма: <strong>{{ .Context.Amount }} рублей</strong></p>
		<a href='{{ .Context.HeroURL }}' target='_blank' rel='noreferrer noopener' class='btn' >
			Войти в личный кабинет
		</a>
		<p>
			Если  появятся вопросы, вот наша почта: {{ .Context.RespondTo }}.
		</p>
		<p>Успехов, <br />команда create.today</p>
	`,
}

var general = Email{
	Subject:  "",
	Template: "default",
	From: EmailSender{
		Email: "hello@createtoday.ru",
		Name:  "CreateToday",
	},
	IsActive: true,
	Type:     "general",
	Context: map[string]interface{}{
		"Domain": "hero.createtoday.ru",
	},
	Body: "",
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
	emails = append(emails, magicLinkLetter, welcomeLetter, orderCreated, general, orderCompleted)
	return &MemoryRepo{}
}
