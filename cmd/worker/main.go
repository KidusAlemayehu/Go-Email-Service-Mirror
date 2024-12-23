package main

import (
	"email-service/config"
	"email-service/internal/dto"
	"email-service/internal/models"
	"email-service/internal/services"
	"email-service/utils/log"
	"encoding/json"
	"os"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const WORKER_COUNT = 10

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

	taskChan := make(chan amqp.Delivery, 100)

	for i := 0; i < WORKER_COUNT; i++ {
		go worker(db, taskChan)
	}

	for d := range msgs {
		taskChan <- d
	}
}

func worker(db *gorm.DB, taskChan <-chan amqp.Delivery) {
	for msg := range taskChan {
		handleTasks(db, msg)
	}
}

func handleTasks(db *gorm.DB, msg amqp.Delivery) {
	log.Logger.Info("Received a message: %s", zap.ByteString("Message Body", msg.Body))

	var task dto.EmailDTO
	if err := json.Unmarshal(msg.Body, &task); err != nil {
		log.Logger.Error("Error Unmarshalling message: %v", zap.Error(err))
		msg.Nack(false, false)
		return
	}

	if err := services.CreateAndSendEmailTask(db, task); err != nil {
		log.Logger.Error("Error sending email: %v", zap.Error(err))
		msg.Nack(false, true)
		return
	} else {
		log.Logger.Info("Successfully processed email for: %s", zap.String("to", task.To))
		msg.Ack(false)
		return
	}
}
