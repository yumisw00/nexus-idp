package queue
import (
	"encoding/json"
	"os"

	"github.com/hibiken/asynq"
)
func NewClient()*asynq.Client{r:=os.Getenv("REDIS_URL");if r==""{r="localhost:6379"};return asynq.NewClient(asynq.RedisClientOpt{Addr:r})}
func EnqueueDocProcess(c *asynq.Client,docID string)error{p,err:=json.Marshal(map[string]string{"doc_id":docID});if err!=nil{return err};_,err=c.Enqueue(asynq.NewTask("job:process_document",p));return err}