package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func createMessage(email string, message string, to string, subject string) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s\r\n", message)

	return msg
}

func SendMail(emailID string, link string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	to := os.Getenv(emailID)
	information := "Click this link to verify your account[Golang Service]: " + link
	subject := "Email Verification from Golang Service"
	msg := []byte(createMessage(emailID, information, to, subject))

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")

	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)

	return err
}
