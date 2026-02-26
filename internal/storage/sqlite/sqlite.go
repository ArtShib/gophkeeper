package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type StoreSqlite struct {
	*sql.DB
}

//go:embed migrations/*.sql
var migrations embed.FS

func New(ctx context.Context, storagePath string) (*StoreSqlite, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("%s: dialect: %w", op, err)
	}

	if err := goose.UpContext(nCtx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &StoreSqlite{DB: db}, nil
}

func (s *StoreSqlite) Stop() error {
	return s.DB.Close()
}
