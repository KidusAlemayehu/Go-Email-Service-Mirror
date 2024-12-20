package services

import (
	"email-service/internal/dto"
	"email-service/utils"
	"fmt"
)

func SendEmailTask(dto dto.EmailDTO) error {

	err := utils.SendMail(dto)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
