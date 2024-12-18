package http

import (
	"email-service/internal/models"
	"encoding/json"
	"net/http"
)

func EmailHandler(rabbitMQ *models.RabbitMQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.EmailTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := rabbitMQ.Publish(task)
		if err != nil {
			http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Email task successfully enqueued"))
	}
}
