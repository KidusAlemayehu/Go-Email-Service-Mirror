package main

import (
	"email-service/config"
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/internal/services"
	"email-service/utils/log"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

func main() {
	config.LoadENV()
	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to initialize database: %v", zap.Error(err))
	}

	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		log.Logger.Error("RABBITMQ_URL environment variable is not set")
		return
	}

	rabbitMQ, err := models.NewRabbitMQ(rmqURL, "email-queue")
	if err != nil {
		log.Logger.Error("Failed to connect to RabbitMQ: %v", zap.Error(err))
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Logger.Error("Failed to consume messages: %v", zap.Error(err))
	}

	for d := range msgs {
		log.Logger.Info("Received a message: %s", zap.ByteString("Message Body", d.Body))

		var task dto.EmailDTO
		if err := json.Unmarshal(d.Body, &task); err != nil {
			log.Logger.Error("Error Unmarshalling message: %v", zap.Error(err))
			d.Nack(false, false)
			continue
		}

		if err := services.CreateAndSendEmailTask(db, task); err != nil {
			log.Logger.Error("Error sending email: %v", zap.Error(err))
			d.Nack(false, true)
		} else {
			log.Logger.Info("Successfully processed email for: %s", zap.String("to", task.To))
			d.Ack(false)
		}
	}
}
