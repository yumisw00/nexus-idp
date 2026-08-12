package worker

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nexus-idp/backend/internal/ai"
)

type DocumentPayload struct {
	DocID string `json:"doc_id"`
	JobID string `json:"job_id"`
}

func HandleDocumentProcess(dbPool *pgxpool.Pool) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p DocumentPayload
		json.Unmarshal(t.Payload(), &p)
		log.Printf("🔥 BINGO! Worker memproses Dokumen: %s", p.DocID)

		filePath := "./uploads/" + p.DocID
		log.Printf("🧠 Mengirim %s ke Gemini 3.5 Flash Lite...", filePath)

		res, err := ai.AnalyzeDocument(ctx, filePath)
		if err != nil {
			log.Printf("❌ Gemini Error (Cek API Key/Quota): %v", err)
			return err
		}

		cleanJSON := strings.TrimSpace(res)
		cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
		cleanJSON = strings.TrimPrefix(cleanJSON, "```")
		cleanJSON = strings.TrimSuffix(cleanJSON, "```")
		cleanJSON = strings.TrimSpace(cleanJSON)

		_, err = dbPool.Exec(ctx, "UPDATE analysis_jobs SET status = 'COMPLETED', res = $1::jsonb, completed_at = NOW() WHERE id = $2", cleanJSON, p.JobID)
		if err != nil {
			log.Printf("DB Update Error: %v", err)
		}

		log.Printf("✅ HASIL EKSTRAKSI GEMINI:\n%s\n===================", res)
		return nil
	}
}
