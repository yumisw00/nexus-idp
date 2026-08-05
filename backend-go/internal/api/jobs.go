package api
import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)
type Job struct{ID string`json:"id"`;Type string`json:"type"`;Status string`json:"status"`}
func GetJobsHandler(db *pgxpool.Pool)http.HandlerFunc{return func(w http.ResponseWriter,r *http.Request){rows,err:=db.Query(context.Background(),"SELECT a.id,d.file_name as type,a.status FROM analysis_jobs a JOIN documents d ON a.document_id=d.id ORDER BY a.started_at DESC LIMIT 10");if err!=nil{http.Error(w,err.Error(),500);return};defer rows.Close();var js []Job;for rows.Next(){var j Job;if err:=rows.Scan(&j.ID,&j.Type,&j.Status);err==nil{js=append(js,j)}};w.Header().Set("Content-Type","application/json");json.NewEncoder(w).Encode(js)}}
