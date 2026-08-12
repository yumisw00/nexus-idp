package queue

import (
	"encoding/json"
	"os"

	"github.com/hibiken/asynq"
)

func NewClient() *asynq.Client {
	r := os.Getenv("REDIS_URL")
	if r == "" {
		r = "localhost:6379"
	}
	return asynq.NewClient(asynq.RedisClientOpt{Addr: r})
}

type DocumentPayload struct {
	DocID string `json:"doc_id"`
	JobID string `json:"job_id"`
}

func EnqueueDocProcess(c *asynq.Client, docID, jobID string) error {
	p, err := json.Marshal(DocumentPayload{DocID: docID, JobID: jobID})
	if err != nil {
		return err
	}
	_, err = c.Enqueue(asynq.NewTask("job:process_document", p))
	return err
}
