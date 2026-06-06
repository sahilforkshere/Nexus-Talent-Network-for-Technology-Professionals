package embedding

import (
	"context"
	"database/sql"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func Init(apiKey string) {
	client = openai.NewClient(apiKey)
}

func EmbedText(ctx context.Context, text string) ([]float32, error) {
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}

func StoreJobEmbedding(ctx context.Context, db *sql.DB, jobID string, vec []float32) error {
	// build postgres vector literal: '[0.1,0.2,...]'
	lit := "["
	for i, v := range vec {
		if i > 0 {
			lit += ","
		}
		lit += fmt.Sprintf("%f", v)
	}
	lit += "]"

	_, err := db.ExecContext(ctx,
		`INSERT INTO job_embeddings (job_id, embedding)
		 VALUES ($1, $2::vector)
		 ON CONFLICT (job_id) DO UPDATE SET embedding = EXCLUDED.embedding`,
		jobID, lit,
	)
	return err
}

func SemanticSearch(ctx context.Context, db *sql.DB, vec []float32, limit int) ([]string, error) {
	lit := "["
	for i, v := range vec {
		if i > 0 {
			lit += ","
		}
		lit += fmt.Sprintf("%f", v)
	}
	lit += "]"

	rows, err := db.QueryContext(ctx,
		`SELECT job_id FROM job_embeddings
		 ORDER BY embedding <-> $1::vector
		 LIMIT $2`,
		lit, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
