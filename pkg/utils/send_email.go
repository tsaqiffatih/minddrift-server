package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/tsaqiffatih/minddrift-server/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	Username         string
	VerificationLink string
	MindDriftEmail   string
}

func SendEmailVerification(cfg *config.Config, toEmail, username, verificationLink string) error {
	emailData := EmailData{
		Username:         username,
		VerificationLink: verificationLink,
		MindDriftEmail:   cfg.MindDriftEmail,
	}

	return sendEmail(cfg, toEmail, "Email Verification - MindDrift", emailData)
}

func sendEmail(cfg *config.Config, toEmail, subject string, data EmailData) error {
	tmpl, err := template.New("email").Parse(EmailVerification)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	d.TLSConfig = nil

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Send email asynchronously
	go func() {
		if err := d.DialAndSend(m); err != nil {
			log.Printf("failed to send email to %s: %v", toEmail, err)
		} else {
			log.Printf("email sent successfully to %s", toEmail)
		}
	}()

	return nil
}

func SendEmail(cfg *config.Config, toEmail, subject, body string) error {
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	d.TLSConfig = nil

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Kirim email secara asynchronous
	go func() {
		if err := d.DialAndSend(m); err != nil {
			log.Printf("failed to send email to %s: %v", toEmail, err)
		} else {
			log.Printf("email sent successfully to %s", toEmail)
		}
	}()

	return nil
}

func GenerateEmailBody(templateString string, data EmailData) (string, error) {
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %v", err)
	}

	return body.String(), nil
}

func getSMTPPort() int {
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return 587
	}

	var portNum int
	fmt.Sscanf(port, "%d", &portNum)
	return portNum
}
