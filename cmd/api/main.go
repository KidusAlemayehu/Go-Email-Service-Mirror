package main

import (
	"email-service/config"
	"email-service/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
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

	router := chi.NewRouter()
	router.Post("/send-email", func(w http.ResponseWriter, r *http.Request) {
		var task models.EmailTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = rabbitMQ.Publish(task)
		if err != nil {
			http.Error(w, "Failed to enqueue email task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Email task successfully enqueued"))
	})

	server_port := os.Getenv("BACKEND_PORT")
	if server_port == "" {
		panic("SERVER_PORT is not set")
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", server_port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Printf("Server started on port :> %s \n", server_port)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
