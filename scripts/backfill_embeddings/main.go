package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://nexus:nexus123@localhost:5432/nexus?sslmode=disable"
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client := openai.NewClient(apiKey)

	rows, err := db.QueryContext(context.Background(),
		`SELECT j.job_id, j.title, j.company, j.description
		 FROM jobs j
		 LEFT JOIN job_embeddings e ON j.job_id = e.job_id
		 WHERE e.job_id IS NULL`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type job struct{ id, title, company, desc string }
	var jobs []job
	for rows.Next() {
		var j job
		if err := rows.Scan(&j.id, &j.title, &j.company, &j.desc); err != nil {
			log.Fatal(err)
		}
		jobs = append(jobs, j)
	}

	fmt.Printf("Found %d jobs without embeddings\n", len(jobs))

	for i, j := range jobs {
		text := j.title + " " + j.company + " " + j.desc
		resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.SmallEmbedding3,
		})
		if err != nil {
			fmt.Printf("[%d/%d] FAILED %s: %v\n", i+1, len(jobs), j.title, err)
			continue
		}
		vec := resp.Data[0].Embedding
		lit := "["
		for k, v := range vec {
			if k > 0 {
				lit += ","
			}
			lit += fmt.Sprintf("%f", v)
		}
		lit += "]"

		_, err = db.ExecContext(context.Background(),
			`INSERT INTO job_embeddings (job_id, embedding)
			 VALUES ($1, $2::vector)
			 ON CONFLICT (job_id) DO UPDATE SET embedding = EXCLUDED.embedding`,
			j.id, lit)
		if err != nil {
			fmt.Printf("[%d/%d] STORE FAILED %s: %v\n", i+1, len(jobs), j.title, err)
			continue
		}
		fmt.Printf("[%d/%d] OK: %s @ %s\n", i+1, len(jobs), j.title, j.company)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Done!")
}
