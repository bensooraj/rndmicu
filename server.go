package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bensooraj/rndmicu/data/database"
	"github.com/bensooraj/rndmicu/graph"
	"github.com/bensooraj/rndmicu/graph/generated"
	aplayground "github.com/bensooraj/rndmicu/graph/playground"
	"github.com/bensooraj/rndmicu/s3engine"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	var err error

	// loads values from .env into the system
	if err = godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	serverCtx, serverCtxCancel := context.WithCancel(context.Background())

	dbDisableTLS := os.Getenv("DB_DISABLE_TLS")
	dbCfg := database.Config{
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_PATH"),
	}
	if dbDisableTLS != "" && dbDisableTLS == "true" {
		dbCfg.DisableTLS = true
	}
	log.Printf("dbCfg: %+v", dbCfg)

	db, err := database.Open(dbCfg)
	if err != nil {
		log.Fatalln("Error establishing db connection:", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	errorChannel := make(chan error, 10)
	// Start the S3 upload workers
	s3Engine := s3engine.NewS3Engine(s3engine.NewS3EngineOptions{
		ErrorChannel:           errorChannel,
		Ctx:                    serverCtx,
		UploadChunkSizeInBytes: 10 * 1024 * 1024,
		Concurrency:            2,
		BucketName:             "rndm.icu",
	})
	err = s3Engine.Init()
	if err != nil {
		log.Println("Error initializing the s3 upload engine")
		return
	}
	s3Engine.StartWorkers()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB:         db,
		S3:         s3Engine,
		CdnBaseURL: os.Getenv("AWS_CDN_BASE_URL"),
	}}))

	// Error handler
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errorChannel:
				log.Println("[ERROR]", err)
			}
		}
	}(serverCtx)

	// "github.com/99designs/gqlgen/graphql/playground"
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// aplayground "github.com/bensooraj/rndmicu/graph/playground"
	http.Handle("/", aplayground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		serverErrors <- http.ListenAndServe(":"+port, nil)
	}()

	// Blocking main and waiting for shutdown of the daemon.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a server error.
	select {
	case e := <-serverErrors:
		log.Printf("server failed to start: %+v", e)
		return

	case sig := <-osSignals:
		// Graceful Shutdown
		log.Println("Received os signal: ", sig.String())
	}

	// Graceful shutdown
	db.Close()
	s3Engine.StopWorkers()
	serverCtxCancel()

	time.Sleep(5 * time.Second) // Allow stuffs to simmer down
}
