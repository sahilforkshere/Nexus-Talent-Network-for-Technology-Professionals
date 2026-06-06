package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Post struct {
	PostID    string
	UserID    string
	Content   string
	CreatedAt string
}

func CreatePost(ctx context.Context, db *sql.DB, userID, content string) (*Post, error) {
	row := db.QueryRowContext(ctx, `
		INSERT INTO posts (author_id, content)
		VALUES ($1, $2)
		RETURNING post_id, author_id, content, created_at::text
	`, userID, content)

	p := &Post{}
	err := row.Scan(&p.PostID, &p.UserID, &p.Content, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPostsByIDs(ctx context.Context, db *sql.DB, ids []string) ([]*Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// filter out job: prefixed items — those are job events, not posts
	var postIDs []string
	for _, id := range ids {
		if len(id) > 4 && id[:4] == "job:" {
			continue
		}
		postIDs = append(postIDs, id)
	}
	if len(postIDs) == 0 {
		return nil, nil
	}

	// build $1,$2,$3... placeholders correctly
	placeholders := ""
	args := make([]any, len(postIDs))
	for i, id := range postIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	rows, err := db.QueryContext(ctx,
		"SELECT post_id, author_id, content, created_at::text FROM posts WHERE post_id IN ("+placeholders+")",
		args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postMap := map[string]*Post{}
	for rows.Next() {
		p := &Post{}
		if err := rows.Scan(&p.PostID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		postMap[p.PostID] = p
	}

	// return in same order as ids
	var posts []*Post
	for _, id := range ids {
		if p, ok := postMap[id]; ok {
			posts = append(posts, p)
		}
	}
	return posts, nil
}
