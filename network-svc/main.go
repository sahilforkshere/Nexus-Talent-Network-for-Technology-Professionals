package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/network-svc/internal/kafka"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to Neo4j
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

	// Start Kafka consumer in background goroutine
	// Runs forever, processing user_created events one by one
	go kafka.ConsumeUserCreated(ctx, driver)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4002"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("network-svc ok"))
	})

	log.Printf("network-svc listening on :%s", port)
	go http.ListenAndServe(":"+port, nil)

	<-ctx.Done()
	log.Println("network-svc shutting down")
}
