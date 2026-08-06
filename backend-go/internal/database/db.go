package database

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	u := os.Getenv("DATABASE_URL")
	if u == "" {
		u = "postgres://root:root@127.0.0.1:5432/nexus?sslmode=disable"
	}

	var p *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		p, err = pgxpool.New(ctx, u)
		if err == nil {
			err = p.Ping(ctx)
			if err == nil {
				c()
				return p, nil
			}
		}
		c()
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
