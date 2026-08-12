package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nexus-idp/backend/internal/queue"
)

func UploadHandler(db *pgxpool.Pool, qc *asynq.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f, h, err := r.FormFile("document")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		fn := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(h.Filename))
		lp := filepath.Join(".", "uploads", fn)

		dst, err := os.Create(lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var did string
		err = db.QueryRow(
			context.Background(),
			"INSERT INTO documents(file_name, file_size_bytes, mime_type, storage_url) VALUES($1, $2, $3, $4) RETURNING id",
			h.Filename, h.Size, h.Header.Get("Content-Type"), lp,
		).Scan(&did)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := queue.EnqueueDocProcess(qc, did, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "queued",
			"doc_id": did,
		})
	}
}
