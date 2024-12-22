package http

import (
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/utils/log"
	"email-service/utils/responses"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SendEmailHandler(db *gorm.DB, rabbitMQ *models.RabbitMQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task dto.EmailDTO
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			errMsg := fmt.Errorf("Invalid request payload")
			log.Logger.Error("Error occured:> ", zap.Error(err))
			responses.NewErrorResponse(errMsg, http.StatusBadRequest, w)
			return
		}

		err := rabbitMQ.Publish(task)
		if err != nil {
			errMsg := fmt.Errorf("Failed to enqueue task")
			log.Logger.Error("Error occured:>", zap.Error(err))
			responses.NewErrorResponse(errMsg, http.StatusInternalServerError, w)
			return
		}
		log.Logger.Info("Email task enqueued successfully")
		responses.NewSuccessResponse("Email task successfully enqueued", nil, http.StatusAccepted, w)
	}
}
