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

func PublishJobPosted(ctx context.Context, jobID, title, company, location string) {
	payload, _ := json.Marshal(map[string]string{
		"job_id":   jobID,
		"title":    title,
		"company":  company,
		"location": location,
	})
	err := writer.WriteMessages(ctx, kafkago.Message{
		Topic: "job_posted",
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka: failed to publish job_posted for %s: %v", jobID, err)
	} else {
		log.Printf("kafka: published job_posted for %s at %s", title, company)
	}
}
