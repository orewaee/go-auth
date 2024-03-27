package email

import (
	"fmt"
	"github.com/orewaee/go-auth/config"
	"net/smtp"
	"strings"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func (mail Mail) Build() string {
	msg := "MIME-version: 1.0;\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\";\n"

	msg += fmt.Sprintf("From: %s\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\n", mail.Subject)
	msg += "\n" + mail.Body + "\n"

	return msg
}

func SendMail(to []string, subject, content string) error {
	from := config.SmtpUsername

	message := Mail{
		from,
		to,
		subject,
		content,
	}.Build()

	return smtp.SendMail(
		config.SmtpHost+":"+config.SmtpPort,
		smtp.PlainAuth(config.SmtpIdentity, from, config.SmtpPassword, config.SmtpHost),
		from, to, []byte(message),
	)
}
