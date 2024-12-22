package models

import (
	"email-service/internal/dto"
	"email-service/utils/log"
	"encoding/json"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Logger.Error("Failed to connect to RabbitMQ: %v", zap.Error(err))
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Logger.Error("Failed to open channel: %v", zap.Error(err))
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (r *RabbitMQ) Publish(task dto.EmailDTO) error {
	body, err := json.Marshal(task)
	if err != nil {
		log.Logger.Error("Error marshaling message: %v", zap.Error(err))
		return err
	}

	err = r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Logger.Error("Error publishing message: %v", zap.Error(err))
		return err
	}

	log.Logger.Info("Message sent to queue: %s", zap.String("to", task.To))
	return nil
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		log.Logger.Error("Error consuming messages: %v", zap.Error(err))
		return nil, err
	}
	log.Logger.Info("Consumer started consuming messages")
	return msgs, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
