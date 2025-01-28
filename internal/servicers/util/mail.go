package util

import (
	"context"
	"fmt"
	"github.com/mosaic-2/IdeYar-server/internal/config"
	"log"
	"net/smtp"
)

const (
	from       = "no-reply@ideyar-app.ir"
	smtpHost   = "smtp-relay.brevo.com"
	smtpPort   = "587"
	mimeHeader = "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"
)

var (
	username string
	password string
)

func init() {
	username = config.GetMailUsername()
	password = config.GetMailPass()
}

// SendEmail handles the common logic for sending emails.
func SendEmail(to []string, subject, messageBody string) error {
	header := fmt.Sprintf("From: %s\r\nSubject: %s\r\nTo: %s\r\n", from, subject, to[0])
	emailMessage := []byte(header + mimeHeader + "\r\n" + messageBody)

	auth := smtp.PlainAuth("", username, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, emailMessage)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func SendSignUpEmail(email, signUpToken, code string) {
	link := fmt.Sprintf("https://ideyar-app.ir/code-verification/%s/%s", signUpToken, code)
	messageBody, err := LoadVerificationEmail(link)
	if err != nil {
		log.Println("Error creating email message:", err)
		return
	}

	err = SendEmail([]string{email}, "Email Verification", messageBody)
	if err != nil {
		log.Println(err)
	}
}

func SendForgetPasswordEmail(ctx context.Context, email, token string) {
	link := fmt.Sprintf("https://ideyar-app.ir/forget-password/%s", token)
	messageBody, err := LoadForgetPasswordEmail(link)
	if err != nil {
		log.Println("Error creating email message:", err)
		return
	}

	err = SendEmail([]string{email}, "Password Reset", messageBody)
	if err != nil {
		log.Println(err)
	}
}

func SendChangeMailEmail(ctx context.Context, email, token string) {
	link := fmt.Sprintf("https://ideyar-app.ir/change-email/%s", token)
	messageBody, err := LoadChangeMailEmail(link)
	if err != nil {
		log.Println("Error creating email message:", err)
		return
	}

	err = SendEmail([]string{email}, "Email Change Verification", messageBody)
	if err != nil {
		log.Println(err)
	}
}
