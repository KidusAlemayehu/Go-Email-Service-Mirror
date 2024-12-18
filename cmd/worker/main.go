package main

import (
	"email-service/config"
	"email-service/internal/models"
	"email-service/internal/services"
	"encoding/json"
	"log"
	"os"
)

func main() {
	config.LoadENV()

	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		log.Println("RABBITMQ_URL environment variable is not set")
		return
	}

	rabbitMQ, err := models.NewRabbitMQ(rmqURL, "email-queue")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		var task models.EmailTask
		if err := json.Unmarshal(d.Body, &task); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			d.Nack(false, false)
			continue
		}

		if err := services.SendEmail(task); err != nil {
			log.Printf("Error sending email: %v", err)
			d.Nack(false, true)
		} else {
			log.Printf("Successfully processed email for: %s", task.To)
			d.Ack(false)
		}
	}
}
