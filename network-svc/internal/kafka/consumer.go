package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	kafkago "github.com/segmentio/kafka-go"
	neo4jdb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/neo4j"
)

// UserCreatedEvent is the message shape published by profile-svc on registration
type UserCreatedEvent struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// ConsumeUserCreated listens to the user_created Kafka topic.
// For every new user registered, creates a Person node in Neo4j.
// Runs forever in a goroutine — one message at a time, in order.
func ConsumeUserCreated(ctx context.Context, driver neo4j.DriverWithContext) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "user_created",
		GroupID:        "network-svc",  // consumer group = Kafka tracks our position
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	log.Println("kafka: listening on topic user_created")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // context cancelled, shutting down
			}
			log.Printf("kafka read error: %v", err)
			continue
		}

		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka: bad message: %v", err)
			continue
		}

		if err := neo4jdb.CreatePersonNode(ctx, driver, event.UserID, event.Name, event.Location); err != nil {
			log.Printf("neo4j: failed to create person node for %s: %v", event.UserID, err)
			continue
		}

		log.Printf("neo4j: created Person node for %s (%s)", event.Name, event.UserID)
	}
}
