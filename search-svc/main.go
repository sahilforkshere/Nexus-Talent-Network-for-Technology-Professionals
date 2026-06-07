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
