package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/graph"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/profile-svc/internal/kafka"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL
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

	// Connect to Neo4j
	neo4jURL := os.Getenv("NEO4J_URL")
	if neo4jURL == "" {
		neo4jURL = "bolt://localhost:7687"
	}
	neo4jDriver, err := neo4j.NewDriverWithContext(neo4jURL, neo4j.BasicAuth("neo4j", "nexuspassword", ""))
	if err != nil {
		log.Fatalf("failed to create neo4j driver: %v", err)
	}
	defer neo4jDriver.Close(context.Background())
	if err = neo4jDriver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("failed to connect to neo4j: %v", err)
	}
	log.Println("connected to neo4j")

	// Init Kafka producer so Register can publish user_created events
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafka.Init(kafkaBroker)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: db, Neo4j: neo4jDriver},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	http.Handle("/", playground.Handler("Profile Service", "/query"))
	http.Handle("/query", auth.Middleware(srv))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("profile-svc ok"))
	})

	log.Printf("profile-svc listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
