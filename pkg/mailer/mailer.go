package mailer

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

const (
	smtpGmailHost   = "smtp.gmail.com"
	smtpGmailPort   = 587
	smtpMailHogHost = "localhost"
	smtpMailHogPort = 1025
)

type EmailSender interface {
	Send(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type emailSender struct {
	dialer *gomail.Dialer
}

func NewEmailSender(username, password, env string) EmailSender {
	if env == "development" {
		return &emailSender{
			dialer: &gomail.Dialer{
				Host: smtpMailHogHost,
				Port: smtpMailHogPort,
			},
		}
	}

	dialer := gomail.NewDialer(smtpGmailHost, smtpGmailPort, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &emailSender{
		dialer: dialer,
	}
}

func (g *emailSender) Send(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	mail := gomail.NewMessage()
	mail.SetHeaders(
		map[string][]string{
			"From":    {"email.admin@eshop.com"},
			"To":      to,
			"Subject": {subject},
		},
	)
	mail.SetHeader("Cc", cc...)
	mail.SetHeader("Bcc", bcc...)
	mail.SetBody("text/html", content)
	if err := g.dialer.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
