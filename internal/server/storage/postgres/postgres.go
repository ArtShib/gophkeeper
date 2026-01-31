package postgres

import (
	"context"
	"database/sql"
	"embed"
)

type StorePostgres struct {
	*sql.DB
}

//go:embed migrations/*.sql
var migrations embed.FS

func New(ctx context.Context, connectionString string) (*StorePostgres, error) {
	return nil, nil
}
