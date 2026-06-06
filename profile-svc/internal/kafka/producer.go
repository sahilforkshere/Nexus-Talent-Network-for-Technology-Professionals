package kafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

var writer *kafkago.Writer

func Init(brokerURL string) {
	writer = &kafkago.Writer{
		Addr:     kafkago.TCP(brokerURL),
		Balancer: &kafkago.LeastBytes{},
	}
}

func PublishUserCreated(ctx context.Context, userID, name, location string) {
	payload, _ := json.Marshal(map[string]string{
		"user_id":  userID,
		"name":     name,
		"location": location,
	})
	err := writer.WriteMessages(ctx, kafkago.Message{
		Topic: "user_created",
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka: failed to publish user_created for %s: %v", userID, err)
	} else {
		log.Printf("kafka: published user_created for %s", name)
	}
}
