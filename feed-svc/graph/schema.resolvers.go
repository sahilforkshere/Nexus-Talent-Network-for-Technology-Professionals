package graph

import (
	"context"
	"fmt"

	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/graph/model"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/cache"
	feeddb "github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/db"
)

// CreatePost saves the post to Postgres and pushes it to the author's own feed in Redis.
func (r *mutationResolver) CreatePost(ctx context.Context, content string) (*model.Post, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	p, err := feeddb.CreatePost(ctx, r.DB, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create post")
	}

	// Push to the author's own feed so they see their own posts
	_ = cache.PushToFeed(ctx, r.Redis, userID, p.PostID)

	return &model.Post{
		PostID:    p.PostID,
		UserID:    p.UserID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}, nil
}

// GetFeed reads the logged-in user's feed from Redis, then fetches full post details from Postgres.
// Redis = fast ordered list of IDs. Postgres = actual content.
func (r *queryResolver) GetFeed(ctx context.Context) ([]*model.FeedItem, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	postIDs, err := cache.GetFeed(ctx, r.Redis, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed")
	}

	posts, err := feeddb.GetPostsByIDs(ctx, r.DB, postIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts")
	}

	var items []*model.FeedItem
	for _, p := range posts {
		items = append(items, &model.FeedItem{
			PostID:    p.PostID,
			UserID:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
		})
	}
	return items, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
