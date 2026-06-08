package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/graph"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/kafka"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	neo4jURL := os.Getenv("NEO4J_URL")
	if neo4jURL == "" {
		neo4jURL = "bolt://localhost:7687"
	}
	driver, err := neo4j.NewDriverWithContext(neo4jURL, neo4j.BasicAuth("neo4j", "nexuspassword", ""))
	if err != nil {
		log.Fatalf("failed to create neo4j driver: %v", err)
	}
	defer driver.Close(ctx)

	if err = driver.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("failed to connect to neo4j: %v", err)
	}
	log.Println("connected to neo4j")

	brokerURL := os.Getenv("KAFKA_BROKER")
	if brokerURL == "" {
		brokerURL = "localhost:9092"
	}
	kafka.InitProducer(brokerURL)
	go kafka.ConsumeUserCreated(ctx, driver)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Neo4j: driver},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4002"
	}

	http.Handle("/", playground.Handler("Network Service", "/query"))
	http.Handle("/query", auth.Middleware(srv))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("network-svc ok"))
	})

	log.Printf("network-svc listening on :%s", port)
	go http.ListenAndServe(":"+port, nil)

	<-ctx.Done()
	log.Println("network-svc shutting down")
}
