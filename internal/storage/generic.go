package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"modernc.org/sqlite"
)

type GenericStorage[T any] struct {
	DB     *sql.DB //*sqlx.DB //*sql.DB
	Table  string
	Driver models.DriverType
}

func NewStorage[T any](db *sql.DB, table string, driverType models.DriverType) *GenericStorage[T] {
	return &GenericStorage[T]{
		DB:     db,
		Table:  table,
		Driver: driverType,
	}
}

func (g *GenericStorage[T]) HandleError(err error, op string) error {

	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	if g.Driver == models.DriverPostgres {
		return g.handlePgError(err, op)
	}

	if g.Driver == models.DriverSQLite {
		return g.handleSqliteError(err, op)
	}

	return fmt.Errorf("%s: %w", op, err)
}

func (g *GenericStorage[T]) handlePgError(err error, op string) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch pgErr.Code {
	case "23505":
		return models.ErrAlreadyExists
	case "23503":
		return models.ErrInvalidReference
	default:
		return fmt.Errorf("%s: postgres[%s]: %w", op, pgErr.Code, err)
	}
}

func (g *GenericStorage[T]) handleSqliteError(err error, op string) error {
	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch sqliteErr.Code() {
	case 2067, 275:
		return models.ErrAlreadyExists
	case 787:
		return models.ErrInvalidReference
	default:
		return fmt.Errorf("%s: sqlite[%d]: %w", op, sqliteErr.Code(), err)
	}
}
