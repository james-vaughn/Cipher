package emailer

import (
	"fmt"
	"net/smtp"
)

type Emailer struct {
	ToAddress      string
	FromAddress    string
	Password       string
	SmtpServerHost string
	SmtpServerPort int
}

func (e *Emailer) Send(subject, body string) error {
	auth := smtp.PlainAuth("", e.FromAddress, e.Password, e.SmtpServerHost)

	smtpServerAddr := fmt.Sprintf("%s:%d", e.SmtpServerHost, e.SmtpServerPort)
	message := e.buildMessage(subject, body)
	err := smtp.SendMail(smtpServerAddr,
		auth,
		e.FromAddress,
		[]string{e.ToAddress},
		[]byte(message))

	if err != nil {
		return fmt.Errorf("Error sending email: \n%v", err)
	}

	return nil
}

func (e *Emailer) buildMessage(subject, body string) string {
	message := ""
	message += fmt.Sprintf("From: %s\n", e.FromAddress)
	message += fmt.Sprintf("To: %s\n", e.ToAddress)
	message += fmt.Sprintf("Subject: %s\n", subject)
	message += "\n" + body

	return message
}
