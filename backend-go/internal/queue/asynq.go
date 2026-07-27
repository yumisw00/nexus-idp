package queue

import (
	"encoding/json"
	"os"

	"github.com/hibiken/asynq"
)

func NewClient() *asynq.Client {
	u := os.Getenv("REDIS_URL")
	if u == "" {
		u = "localhost:6379"
	}
	return asynq.NewClient(asynq.RedisClientOpt{Addr: u})
}

func EnqueueDocProcess(c *asynq.Client, id string) error {
	p, err := json.Marshal(map[string]string{"doc_id": id})
	if err != nil {
		return err
	}
	_, err = c.Enqueue(asynq.NewTask("job:process_document", p))
	return err
}
