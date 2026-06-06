package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const feedSize = 50

// PushToFeed adds a post to a user's feed sorted set.
// Redis sorted set key: feed:{userID}
// Score = unix timestamp so feed is always newest-first when reversed.
func PushToFeed(ctx context.Context, rdb *redis.Client, userID, postID string) error {
	key := "feed:" + userID
	return rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: postID,
	}).Err()
}

// GetFeed returns the latest post IDs for a user (newest first, capped at feedSize).
func GetFeed(ctx context.Context, rdb *redis.Client, userID string) ([]string, error) {
	key := "feed:" + userID
	// ZREVRANGE = highest score first = newest posts first
	return rdb.ZRevRange(ctx, key, 0, int64(feedSize-1)).Result()
}
