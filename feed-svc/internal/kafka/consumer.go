package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/redis/go-redis/v9"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/cache"
)

type JobPostedEvent struct {
	JobID    string `json:"job_id"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
}

// ConsumeJobPosted listens to the job_posted Kafka topic.
// For every new job, it pushes the job_id into every user's feed in Redis.
// This is how "Sahil opens feed and sees a new Google job" works automatically.
func ConsumeJobPosted(ctx context.Context, db *sql.DB, rdb *redis.Client, brokerURL string) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        []string{brokerURL},
		Topic:          "job_posted",
		GroupID:        "feed-svc",
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	log.Println("kafka: listening on topic job_posted")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("kafka read error: %v", err)
			continue
		}

		var event JobPostedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka: bad message: %v", err)
			continue
		}

		// Push job to all users' feeds
		pushed, err := pushJobToAllFeeds(ctx, db, rdb, event)
		if err != nil {
			log.Printf("feed: failed to push job %s to feeds: %v", event.JobID, err)
			continue
		}
		log.Printf("feed: pushed job '%s' (%s) to %d user feeds", event.Title, event.Company, pushed)
	}
}

// pushJobToAllFeeds fetches all user IDs and adds the job to each feed
func pushJobToAllFeeds(ctx context.Context, db *sql.DB, rdb *redis.Client, event JobPostedEvent) (int, error) {
	rows, err := db.QueryContext(ctx, "SELECT user_id FROM users")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		// Use job_id as the feed item — prefixed so getFeed can tell it's a job
		if err := cache.PushToFeed(ctx, rdb, userID, "job:"+event.JobID); err != nil {
			log.Printf("redis: failed to push to feed:%s: %v", userID, err)
			continue
		}
		count++
	}
	return count, nil
}
