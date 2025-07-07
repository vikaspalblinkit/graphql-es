package main

import (
	"log"
	"net/http"
	"os"

	"compensation-api/graph"
	"compensation-api/internal/elastic"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	esv8 "github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := esv8.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}
	esClient := &elastic.Client{Client: es, Index: "compensations"}

	if err := esClient.CreateIndex("index.json"); err != nil {
		log.Printf("Failed to create index: %v", err)
	}
	if err := esClient.BulkUploadCSV("./dataset"); err != nil {
		log.Printf("Failed to bulk upload CSV: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{ES: esClient}}))

	http.Handle("/", playground.Handler("GraphQL", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/compensation_data", compensationDataHandler(esClient))

	log.Printf("🚀 Server started at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
