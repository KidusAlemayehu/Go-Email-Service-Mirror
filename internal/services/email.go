package services

import (
	"bytes"
	"email-service/internal/models"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
)

func SendEmail(task models.EmailTask) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", username, password, smtpHost)

	var emailBody bytes.Buffer
	writer := multipart.NewWriter(&emailBody)

	headers := []string{
		fmt.Sprintf("From: %s", task.From),
		fmt.Sprintf("To: %s", task.To),
		fmt.Sprintf("Subject: %s", task.Subject),
		"Mime-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s", writer.Boundary()),
	}
	if len(task.Cc) > 0 {
		headers = append(headers, fmt.Sprintf("Cc: %s", joinAddresses(task.Cc)))
	}

	for _, h := range headers {
		emailBody.WriteString(h + "\r\n")
	}
	emailBody.WriteString("\r\n")

	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=\"utf-8\""},
	})
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}
	_, err = textPart.Write([]byte(task.Body))
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if task.Attachments != nil {
		for _, attachment := range task.Attachments {
			if attachment.Data == "" {
				continue
			}

			decodedData, err := base64.StdEncoding.DecodeString(attachment.Data)
			if err != nil {
				return fmt.Errorf("failed to decode attachment: %w", err)
			}

			part, err := writer.CreatePart(textproto.MIMEHeader{
				"Content-Type":              []string{attachment.ContentType},
				"Content-Disposition":       []string{fmt.Sprintf(`attachment; filename="%s"`, attachment.Filename)},
				"Content-Transfer-Encoding": []string{"base64"},
			})
			if err != nil {
				return fmt.Errorf("failed to create attachment part: %w", err)
			}
			_, err = part.Write(decodedData)
			if err != nil {
				return fmt.Errorf("failed to write attachment: %w", err)
			}
		}
	}

	writer.Close()

	recipients := []string{task.To}
	if task.Cc != nil {
		recipients = append(recipients, task.Cc...)
	}
	if task.Bcc != nil {
		recipients = append(recipients, task.Bcc...)
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		task.From,
		recipients,
		emailBody.Bytes(),
	)
}

func joinAddresses(addresses []string) string {
	return fmt.Sprintf("%s", bytes.Join([][]byte{[]byte(addresses[0])}, []byte(", ")))
}
