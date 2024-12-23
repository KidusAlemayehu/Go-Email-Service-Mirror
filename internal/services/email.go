package services

import (
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/utils"
	"email-service/utils/log"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const maxRetries = 5

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
		log.Logger.Error("Error creating email task: %v", zap.Error(err))
		return fmt.Errorf("failed to create email task in DB: %w", err)
	}
	log.Logger.Info("Email task created successfully. Processing: %s", zap.String("to", task.To))
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
			log.Logger.Info("Email sent successfully. Task ID: %s", zap.String("id", emailTask.ID.String()))
			return nil
		}

		// If error occurs, handle retry and update task status
		emailTask.Status = "failed"
		emailTask.RetryCount = retries + 1
		emailTask.LastAttempt = time.Now()
		db.Save(&emailTask)
		backoff := time.Duration(math.Pow(2, float64(retries))) * time.Second
		log.Logger.Error("Error sending email. Retrying. Task ID: %s, Error: %v", zap.String("id", emailTask.ID.String()), zap.Error(err))
		time.Sleep(backoff)
	}

	// If maxRetries reached and the email failed, log it and stop retrying
	emailTask.Status = "failed"
	emailTask.RetryCount = maxRetries
	emailTask.LastAttempt = time.Now()
	db.Save(&emailTask)
	log.Logger.Error(fmt.Sprintf("Email sending failed after %d retries.", maxRetries), zap.String("task_id", emailTask.ID.String()), zap.Error(fmt.Errorf("Connection Timeout, Maximum retry limit reached.")))

	return fmt.Errorf("email sending failed after %d retries", maxRetries)
}
