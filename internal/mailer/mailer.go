package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/wneessen/go-mail"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(cfg config.Config, logger *jsonlog.Logger) Mailer {
	client, err := mail.NewClient(
		cfg.SMTP.Host,
		mail.WithPort(cfg.SMTP.Port),
		mail.WithUsername(cfg.SMTP.Username),
		mail.WithPassword(cfg.SMTP.Password),
		mail.WithTimeout(5*time.Second),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
	)

	if err != nil {
		logger.LogFatal(err, nil)
	}

	logger.LogInfo("mail client created", nil)

	return Mailer{
		client: client,
		sender: cfg.SMTP.Sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return nil
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	msg.To(recipient)
	msg.From(m.sender)
	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	err = m.client.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
