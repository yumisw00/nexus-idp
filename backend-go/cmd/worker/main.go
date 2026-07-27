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
		log.Fatalf("DB init failed: %v", err)
	}
	defer dbPool.Close()

	u := os.Getenv("REDIS_URL")
	if u == "" {
		u = "localhost:6379"
	}

	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: u}, asynq.Config{})
	mux := asynq.NewServeMux()
	mux.HandleFunc("job:process_document", worker.HandleDocumentProcess(dbPool))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Worker start failed: %v", err)
	}
}
