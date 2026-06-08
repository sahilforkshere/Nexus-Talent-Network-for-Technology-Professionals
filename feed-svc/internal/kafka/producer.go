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

func PublishPostCreated(ctx context.Context, postID, userID, content string) {
	payload, _ := json.Marshal(map[string]string{
		"post_id": postID,
		"user_id": userID,
		"content": content,
	})
	err := writer.WriteMessages(ctx, kafkago.Message{
		Topic: "post_created",
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka: failed to publish post_created %s: %v", postID, err)
	} else {
		log.Printf("kafka: published post_created %s by %s", postID, userID)
	}
}
