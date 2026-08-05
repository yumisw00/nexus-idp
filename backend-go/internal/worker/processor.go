package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/nexus-idp/backend/internal/ai"
)

type DocumentPayload struct {
	DocumentID string `json:"doc_id"`
}

func HandleDocumentProcess(db *pgxpool.Pool) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload DocumentPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return err
		}

		var storageURL string
		err := db.QueryRow(ctx, "SELECT storage_url FROM documents WHERE id=$1", payload.DocumentID).Scan(&storageURL)
		if err != nil {
			return err
		}

		var jobID string
		err = db.QueryRow(ctx, "INSERT INTO analysis_jobs(document_id, status, started_at) VALUES($1, 'PROCESSING', NOW()) RETURNING id", payload.DocumentID).Scan(&jobID)
		if err != nil {
			return err
		}

		response, analyzeErr := ai.AnalyzeDocument(ctx, storageURL)
		
		status := "COMPLETED"
		if analyzeErr != nil {
			status = "FAILED"
		}

		_, updateErr := db.Exec(ctx, "UPDATE analysis_jobs SET status=$1, raw_ai_response=$2, completed_at=NOW() WHERE id=$3", status, response, jobID)
		if updateErr != nil {
			return updateErr
		}

		return analyzeErr
	}
}