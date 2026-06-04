package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	svcName := "feed-svc"
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s ok", svcName)
	})
	log.Printf("%s listening on :%s", svcName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
