package config

import "gopkg.in/gomail.v2"

const CONFIG_SMTP_HOST = "smtp.mailtrap.io"
const CONFIG_SMTP_PORT = 2525
const CONFIG_SENDER_NAME = "Fadhli dev <emailanda@gmail.com>"
const CONFIG_AUTH_EMAIL = "2ea2d9a2eeebf2"
const CONFIG_AUTH_PASSWORD = "bb9603f59a99a1"

func SendEmail(recipient, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", recipient)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
