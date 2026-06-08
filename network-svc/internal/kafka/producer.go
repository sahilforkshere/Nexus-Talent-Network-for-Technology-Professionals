package kafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

var writer *kafkago.Writer

func InitProducer(brokerURL string) {
	writer = &kafkago.Writer{
		Addr:     kafkago.TCP(brokerURL),
		Balancer: &kafkago.LeastBytes{},
	}
}

func PublishConnectionAccepted(ctx context.Context, fromUserID, toUserID string) {
	payload, _ := json.Marshal(map[string]string{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	})
	err := writer.WriteMessages(ctx, kafkago.Message{
		Topic: "connection_accepted",
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka: failed to publish connection_accepted (%s↔%s): %v", fromUserID, toUserID, err)
	} else {
		log.Printf("kafka: published connection_accepted (%s↔%s)", fromUserID, toUserID)
	}
}
