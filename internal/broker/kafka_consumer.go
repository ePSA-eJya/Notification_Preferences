package broker

import (
	"context"
	"encoding/json"
	"log"

	"Notification_Preferences/internal/entities"
	notificationUseCase "Notification_Preferences/internal/notification/usecase"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader       *kafka.Reader
	notifUseCase notificationUseCase.NotificationService
}

func NewKafkaConsumer(brokerURLs []string, topic string, groupID string, notifUseCase notificationUseCase.NotificationService) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokerURLs,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader:       r,
		notifUseCase: notifUseCase,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka Consumer for topic: %s", c.reader.Config().Topic)
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping consumer...")
			if err := c.reader.Close(); err != nil {
				log.Printf("Error closing reader: %v", err)
			}
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			var event map[string]interface{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event map: %v", err)
				continue
			}

			// We need to parse this into entities.Event carefully
			// Let's just unmarshal into the strict struct directly first
			var typedEvent entities.Event
			if err := json.Unmarshal(msg.Value, &typedEvent); err != nil {
				log.Printf("Error unmarshaling strictly: %v", err)
			} else {
				log.Printf("Received Event Action: %s", typedEvent.ActionType)
				if err := c.notifUseCase.ProcessEvent(ctx, &typedEvent); err != nil {
					log.Printf("Error processing event: %v", err)
				}
			}
		}
	}
}
