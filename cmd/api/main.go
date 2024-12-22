package main

import (
	"email-service/config"
	httpHandlers "email-service/internal/http"
	"email-service/internal/models"
	"email-service/utils/log"
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	config.LoadENV()

	db, err := config.InitDB()
	if err != nil {
		log.Logger.Info("Database initialization failed", zap.Error(err))
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

	router := chi.NewRouter()
	router.Post("/send-email", httpHandlers.SendEmailHandler(db, rabbitMQ))

	server_port := os.Getenv("BACKEND_PORT")
	if server_port == "" {
		log.Logger.Error("Variable SERVER_PORT is not set")
		panic("SERVER_PORT is not set")
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", server_port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Logger.Info("Server started on port :> %s \n", zap.String("SERVER_PORT", server_port))
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
