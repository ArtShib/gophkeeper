package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type StorePostgres struct {
	*sql.DB
}

//go:embed migrations/*.sql
var migrations embed.FS

func New(ctx context.Context, connectionString string) (*StorePostgres, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	goose.SetBaseFS(migrations)
	if err := goose.UpContext(nCtx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &StorePostgres{DB: db}, nil
}
