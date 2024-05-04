package hero

import (
	"bytes"
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/config"
	"createtodayapi/internal/entity"
	"createtodayapi/internal/logger"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"os"
)

const pathToTemplates = "/internal/emails/"
const templateExtension = ".html"

type EmailsService struct {
	config *config.Config
	repo   *MemoryRepo
}

type IEmailsService interface {
	GetEmailByType(context context.Context, emailType string) (*entity.Email, error)
	SendEmail(email *entity.Email, to []string) error
}

func (s *EmailsService) GetEmailByType(context context.Context, emailType string) (*entity.Email, error) {
	email, err := s.repo.FindByType(context, emailType)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, common.ErrInternalError
	}
	return email, nil
}

func (s *EmailsService) buildTemplatePath(templateName string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd + pathToTemplates + templateName + templateExtension, nil
}

func (s *EmailsService) buildEmailHtml(email *entity.Email) (string, error) {

	body, err := template.New("").Parse(email.Body)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not parse html body in template %s ", email.Template), "error", err)
		return "", err
	}

	var bodyBuffer bytes.Buffer

	err = body.Execute(&bodyBuffer, email)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not execute html body in template %s ", email.Template), "error", err)
		return "", err
	}

	templatePath, err := s.buildTemplatePath(email.Template)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not build template path for %s ", email.Template), "error", err)
		return "", err
	}

	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not parse template %s ", email.Template), "error", err)
		return "", err
	}

	var tpl bytes.Buffer

	email.Context["Body"] = template.HTML(bodyBuffer.String())

	err = templ.Execute(&tpl, email)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not execute template %s ", email.Template), "error", err)
	}

	return tpl.String(), nil
}

func (s *EmailsService) buildEmailSenderName(sender entity.EmailSender) string {
	return sender.Name + " " + "<" + sender.Email + ">"
}

func (s *EmailsService) SendEmail(email *entity.Email, to []string) error {
	htmlBody, err := s.buildEmailHtml(email)

	if err != nil {
		return err
	}

	charSet := "utf-8"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.config.AwsRegion)},
	)

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(to),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(s.buildEmailSenderName(email.From)),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		return err
	}

	logger.Log.Info(fmt.Sprintf("sent email with id %v", result.MessageId))

	return nil
}

func NewEmailService(config *config.Config, repo *MemoryRepo) *EmailsService {
	return &EmailsService{
		config: config,
		repo:   repo,
	}
}
