package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bensooraj/rndmicu/data/database"
	"github.com/bensooraj/rndmicu/graph"
	"github.com/bensooraj/rndmicu/graph/generated"
	aplayground "github.com/bensooraj/rndmicu/graph/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
}

func main() {
	dbCfg := database.Config{
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_PATH"),
	}
	log.Println("dbCfg: ", dbCfg)

	db, err := database.Open(dbCfg)
	if err != nil {
		log.Fatalln("Error establishing db connection:", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	// "github.com/99designs/gqlgen/graphql/playground"
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// aplayground "github.com/bensooraj/rndmicu/graph/playground"
	http.Handle("/", aplayground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
