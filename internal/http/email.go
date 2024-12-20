package http

import (
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/utils/responses"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendEmailHandler(rabbitMQ *models.RabbitMQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task dto.EmailDTO
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			errMsg := fmt.Errorf("Invalid request payload")
			log.Fatal(err)
			responses.NewErrorResponse(errMsg, http.StatusBadRequest, w)
			return
		}

		err := rabbitMQ.Publish(task)
		if err != nil {
			errMsg := fmt.Errorf("Failed to enqueue task")
			log.Fatal(errMsg)
			responses.NewErrorResponse(errMsg, http.StatusInternalServerError, w)
			return
		}
		responses.NewSuccessResponse("Email task successfully enqueued", nil, http.StatusAccepted, w)
	}
}
