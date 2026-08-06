package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/nexus-idp/backend/internal/ai"
)

type DocumentPayload struct {
	DocID string `json:"doc_id"`
}

func HandleDocumentProcess(ctx context.Context, t *asynq.Task) error {
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
	
	log.Printf("✅ HASIL EKSTRAKSI GEMINI:\n%s\n===================", res)
	return nil
}
