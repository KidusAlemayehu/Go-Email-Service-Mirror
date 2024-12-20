package models

import (
	"time"

	"github.com/google/uuid"
)

// EmailTask represents an email task that can be stored in the database.
type EmailTask struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	To          string       `gorm:"not null" json:"to"`
	From        string       `gorm:"not null" json:"from"`
	Subject     string       `gorm:"not null" json:"subject"`
	Body        string       `gorm:"not null" json:"body"`
	Cc          []string     `gorm:"type:text[]" json:"cc,omitempty"`
	Bcc         []string     `gorm:"type:text[]" json:"bcc,omitempty"`
	ReplyTo     string       `gorm:"type:varchar(255);default:null" json:"reply_to,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:EmailTaskID" json:"attachments"`
	Status      string       `gorm:"not null;default:'pending'" json:"status"`
	RetryCount  int          `gorm:"default:0" json:"retry_count"`
	LastAttempt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Attachment represents an email attachment linked to an EmailTask.
type Attachment struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EmailTaskID uuid.UUID `gorm:"type:uuid;not null" json:"email_task_id"`
	Filename    string    `gorm:"not null" json:"filename"`
	ContentType string    `gorm:"not null" json:"content_type"`
	Data        string    `gorm:"not null" json:"data"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
