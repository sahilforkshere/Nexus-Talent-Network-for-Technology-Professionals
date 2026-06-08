package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/search-svc/graph"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/search-svc/internal/auth"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/search-svc/internal/proximity"
	"github.com/sahilpal/Nexus-TalentNetworkForTechnologyProfessionals/search-svc/internal/search"
)

func main() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}
	if err := search.Init(esURL); err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}
	log.Println("connected to elasticsearch")

	neo4jURI := os.Getenv("NEO4J_URI")
	if neo4jURI == "" {
		neo4jURI = "bolt://localhost:7687"
	}
	neo4jUser := os.Getenv("NEO4J_USER")
	if neo4jUser == "" {
		neo4jUser = "neo4j"
	}
	neo4jPass := os.Getenv("NEO4J_PASSWORD")
	if neo4jPass == "" {
		neo4jPass = "nexuspassword"
	}
	if err := proximity.Init(neo4jURI, neo4jUser, neo4jPass); err != nil {
		log.Printf("warning: neo4j unavailable, proximity boost disabled: %v", err)
	} else {
		log.Println("connected to neo4j")
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4005"
	}

	http.Handle("/", playground.Handler("Search Service", "/query"))
	http.Handle("/query", auth.Middleware(srv))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("search-svc ok"))
	})

	log.Printf("search-svc listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
