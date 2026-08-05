package main

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	
	"github.com/nexus-idp/backend/internal/database"
	"github.com/nexus-idp/backend/internal/worker"
)

func main() {
	dbPool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisURL},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("job:process_document", worker.HandleDocumentProcess(dbPool))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("worker failed: %v", err)
	}
}