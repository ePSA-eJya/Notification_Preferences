package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer establishes connection to localhost:9092
func NewKafkaProducer(brokerURL string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: w}
}

// Publish event takes any data and publishes it to kafka
func (p *KafkaProducer) Publish(ctx context.Context, topic string, payload interface{}) error {
	log.Printf("[KafkaProducer] publish requested topic=%s payloadType=%T", topic, payload)

	// 1. Convert the Go Map/Struct into standard JSON format
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal event data: %v", err)
		return err
	}
	log.Printf("[KafkaProducer] marshaled payload size=%d bytes", len(jsonData))

	// 2. Create the Kafka Message
	msg := kafka.Message{
		Value: jsonData,
	}

	// 3. Send the message to the Broker!
	log.Printf("[KafkaProducer] writing message to kafka")
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	log.Printf("[KafkaProducer] publish succeeded topic=%s", topic)

	return nil
}
