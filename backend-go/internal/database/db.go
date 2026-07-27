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
		u = "postgres://root:root@localhost:5432/nexus?sslmode=disable"
	}
	p, err := pgxpool.New(context.Background(), u)
	if err != nil {
		return nil, err
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.Ping(c); err != nil {
		return nil, err
	}
	return p, nil
}
