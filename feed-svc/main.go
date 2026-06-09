package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/graph"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/feed-svc/internal/kafka"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}
	rdb := redis.NewClient(opt)
	log.Println("connected to redis")

	brokerURL := os.Getenv("KAFKA_BROKER")
	if brokerURL == "" {
		brokerURL = "localhost:9092"
	}
	kafka.InitProducer(brokerURL)
	go kafka.ConsumeJobPosted(ctx, db, rdb, brokerURL)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: db, Redis: rdb},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(200))

	port := os.Getenv("PORT")
	if port == "" {
		port = "4004"
	}

	http.Handle("/", playground.Handler("Feed Service", "/query"))
	http.Handle("/query", auth.Middleware(srv))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("feed-svc ok"))
	})

	log.Printf("feed-svc listening on :%s", port)
	go http.ListenAndServe(":"+port, nil)

	<-ctx.Done()
	log.Println("feed-svc shutting down")
}
