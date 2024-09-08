package email

import (
	"gopkg.in/gomail.v2"
)

func SendMail(from string, to string, subject string, msg string, envMap map[string]string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, envMap["MAIL_USERNAME"], envMap["MAIL_PASSWORD"])

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
