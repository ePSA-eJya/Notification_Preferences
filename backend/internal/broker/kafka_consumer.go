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
		MinBytes: 1,    // 10KB
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
			c.processNextMessage(ctx)
		}
	}
}

func (c *KafkaConsumer) processNextMessage(ctx context.Context) {
	// Recover from panics so the consumer goroutine doesn't die silently
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔴 PANIC in Kafka consumer (recovered): %v", r)
		}
	}()

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		log.Printf("Error reading message: %v", err)
		return
	}

	log.Printf("📩 Raw Kafka message: %s", string(msg.Value))

	var typedEvent entities.Event
	if err := json.Unmarshal(msg.Value, &typedEvent); err != nil {
		log.Printf("Error unmarshaling event: %v", err)
		return
	}

	log.Printf("✅ Parsed Event — Action: %s, ActorID: %s, EntityID: %s, EntityType: %s",
		typedEvent.ActionType, typedEvent.ActorID, typedEvent.EntityID, typedEvent.EntityType)

	if err := c.notifUseCase.ProcessEvent(ctx, &typedEvent); err != nil {
		log.Printf("❌ Error processing event: %v", err)
	} else {
		log.Printf("✅ Successfully processed %s event", typedEvent.ActionType)
	}
}
