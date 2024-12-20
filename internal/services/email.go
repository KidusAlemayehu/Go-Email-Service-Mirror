package services

import (
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const maxRetries = 5
const retryDelay = 5 * time.Second

// CreateAndSendEmailTask creates the initial email task and attempts to send it
func CreateAndSendEmailTask(db *gorm.DB, task dto.EmailDTO) error {
	var attachments []models.Attachment
	for _, dtoAttachment := range task.Attachments {
		attachments = append(attachments, models.Attachment{
			Filename:    dtoAttachment.Filename,
			ContentType: dtoAttachment.ContentType,
			Data:        dtoAttachment.Data,
		})
	}
	emailTask := models.EmailTask{
		From:        task.From,
		To:          task.To,
		Cc:          task.Cc,
		Bcc:         task.Bcc,
		Subject:     task.Subject,
		Body:        task.Body,
		ReplyTo:     task.ReplyTo,
		Attachments: attachments,
		Status:      "pending",
	}

	if err := db.Create(&emailTask).Error; err != nil {
		return fmt.Errorf("failed to create email task in DB: %w", err)
	}

	return processEmailTask(db, task, &emailTask)
}

// processEmailTask attempts to send the email and handles retries.
func processEmailTask(db *gorm.DB, task dto.EmailDTO, emailTask *models.EmailTask) error {
	// Try sending the email task up to maxRetries
	for retries := 0; retries < maxRetries; retries++ {
		err := utils.SendMail(task)
		if err == nil {
			// If email is sent successfully, update the task status to "sent"
			emailTask.Status = "sent"
			emailTask.RetryCount = retries
			emailTask.LastAttempt = time.Now()
			db.Save(&emailTask)
			return nil
		}

		// If error occurs, handle retry and update task status
		emailTask.Status = "failed"
		emailTask.RetryCount = retries + 1
		emailTask.LastAttempt = time.Now()
		db.Save(&emailTask)

		// Sleep before retrying
		time.Sleep(retryDelay)
	}

	// If maxRetries reached and the email failed, log it and stop retrying
	emailTask.Status = "failed"
	emailTask.RetryCount = maxRetries
	emailTask.LastAttempt = time.Now()
	db.Save(&emailTask)
	return fmt.Errorf("email sending failed after %d retries", maxRetries)
}
