package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/graph"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/embedding"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/kafka"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/jobs-svc/internal/search"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://nexus:nexus@localhost:5432/nexus?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("connected to postgres")

	// ensure pgvector extension and job_embeddings table exist
	_, _ = db.Exec(`CREATE EXTENSION IF NOT EXISTS vector`)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS job_embeddings (
		job_id UUID PRIMARY KEY REFERENCES jobs(job_id) ON DELETE CASCADE,
		embedding vector(1536) NOT NULL
	)`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS job_embeddings_hnsw ON job_embeddings USING hnsw (embedding vector_l2_ops)`)

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}
	if err := search.Init(esURL); err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}
	log.Println("connected to elasticsearch")

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafka.Init(kafkaBroker)

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY env var is required")
	}
	embedding.Init(openAIKey)
	log.Println("embedding client initialized")

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: db},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4003"
	}

	http.Handle("/", playground.Handler("Jobs Service", "/query"))
	http.Handle("/query", auth.Middleware(srv))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("jobs-svc ok"))
	})

	log.Printf("jobs-svc listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
