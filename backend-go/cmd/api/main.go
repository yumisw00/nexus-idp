package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	
	"github.com/nexus-idp/backend/internal/database"
	"github.com/nexus-idp/backend/internal/queue"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

func main() {
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("failed to create uploads directory: %v", err)
	}

	dbPool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	queueClient := queue.NewClient()
	defer queueClient.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json"); w.WriteHeader(http.StatusOK); w.Write([]byte(`{"status":"ok"}`));
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/jobs", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"id": "SYSTEM_ONLINE", "type": "Engine", "status": "Ready", "progress": 100}]`))
		})
		r.Post("/documents", func(w http.ResponseWriter, req *http.Request) {
			log.Println("DEBUG: Handler invoked")
			if err := req.ParseMultipartForm(10 << 20); err != nil {
				log.Printf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			file, handler, err := req.FormFile("document")
			if err != nil {
				log.Printf("Error getting file: %v", err)
				http.Error(w, "File missing", http.StatusBadRequest)
				return
			}
			defer file.Close()
			err = queue.EnqueueDocProcess(queueClient, handler.Filename)
			if err != nil { log.Printf("Queue error: %v", err) }
			log.Printf("QUEUED: %s", handler.Filename)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}