package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nexus-idp/backend/internal/ai"
)

func HandleDocumentProcess(db *pgxpool.Pool) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p map[string]string
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		id := p["doc_id"]
		var url string
		if err := db.QueryRow(ctx, "SELECT url FROM documents WHERE id = $1", id).Scan(&url); err != nil {
			return err
		}
		if _, err := db.Exec(ctx, "UPDATE documents SET status = 'PROCESSING' WHERE id = $1", id); err != nil {
			return err
		}
		res, err := ai.AnalyzeDocument(ctx, url)
		status := "COMPLETED"
		if err != nil {
			status = "FAILED"
		}
		if _, dbErr := db.Exec(ctx, "UPDATE documents SET status = $1 WHERE id = $2", status, id); dbErr != nil {
			return dbErr
		}
		if err == nil {
			if _, dbErr := db.Exec(ctx, "INSERT INTO analysis_jobs (document_id, result) VALUES ($1, $2)", id, res); dbErr != nil {
				return dbErr
			}
		}
		return err
	}
}
