package utils

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(toEmail, token string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	baseURL := emailVerificationUrl

	verificationLink := fmt.Sprintf("%s?token=%s", baseURL, token)
	log.Println(verificationLink)

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Click the link below to verify your email:</p>
		<a href="%s">%s</a>
	`, verificationLink, verificationLink))

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, emailPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
