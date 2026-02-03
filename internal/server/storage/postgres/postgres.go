package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/ArtShib/gophkeeper/internal/server/models"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type StorePostgres struct {
	db *sql.DB
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
	return &StorePostgres{db: db}, nil
}

func (pg *StorePostgres) Close() error {
	return pg.db.Close()
}

// AddUser регистрация пользователя
func (pg *StorePostgres) AddUser(ctx context.Context, login string, passHash []byte, createdAT int64) (*models.User, error) {
	const op = "storage.postgres.AddUser"
	var user models.User
	query := `
			INSERT INTO gophkeeper.users (login, pass_hash, created_at) 
			VALUES ($1, $2, $3) RETURNING id, login, pass_hash`

	if err := pg.db.QueryRowContext(ctx, query, login, passHash, createdAT).Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserExists)
			}
		}
		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// GetUser select пользователя по логину
func (pg *StorePostgres) GetUser(ctx context.Context, login string) (*models.User, error) {
	const op = "storage.postgres.GetUser"

	query := "SELECT id, login, pass_hash FROM gophkeeper.users WHERE login = $1"

	row := pg.db.QueryRowContext(ctx, query, login)

	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// AddSecret добавление секрета
func (pg *StorePostgres) AddSecret(
	ctx context.Context,
	secretID string,
	ownerID int64,
	typeSecret string,
	data []byte,
	metadata string,
	createdAT int64) error {
	const op = "storage.postgres.AddSecret"
	query := `
			INSERT INTO gophkeeper.data (secret_id, owner_id, type, data, metadata, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pg.db.ExecContext(ctx, query, secretID, ownerID, typeSecret, data, metadata, createdAT)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetUserSecrets список секретов
func (pg *StorePostgres) GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error) {
	const op = "storage.postgres.GetUserSecrets"

	query := `
		SELECT secret_id, type, data, metadata, created_at, updated_at
       FROM gophkeeper.data
       WHERE owner_id = $1 AND is_deleted = FALSE`

	rows, err := pg.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var secrets models.ArraySecret

	for rows.Next() {
		var secret models.Secret
		var updatedAt sql.NullInt64

		err := rows.Scan(&secret.ID, &secret.Type, &secret.Data, &secret.Metadata, &secret.CreatedAt, &updatedAt)

		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		if updatedAt.Valid {
			secret.UpdatedAt = updatedAt.Int64
		}

		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iter: %w", op, err)
	}

	return secrets, nil
}

// UpdateSecret обновление секрета
func (pg *StorePostgres) UpdateSecret(
	ctx context.Context,
	secretID string,
	ownerID int64,
	typeSecret string,
	data []byte,
	metadata string,
	updatedAT int64) error {
	const op = "storage.postgres.AddSecret"

	query := `UPDATE gophkeeper.data 
          SET data = $1, metadata = $2, updated_at = $3 
          WHERE secret_id = $4 AND owner_id = $5 AND type = $6`

	_, err := pg.db.ExecContext(ctx, query, data, metadata, updatedAT, secretID, ownerID, typeSecret)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteSecret установка признака удаления
func (pg *StorePostgres) DeleteSecret(
	ctx context.Context,
	secretID string,
	ownerID int64,
	typeSecret string,
	updatedAT int64) error {
	const op = "storage.postgres.AddSecret"
	query := `
			UPDATE gophkeeper.data 
			SET is_deleted = True, updated_at = $1 
			WHERE secret_id = $2 AND owner_id = $3 AND type = $4`

	_, err := pg.db.ExecContext(ctx, query, updatedAT, secretID, ownerID, typeSecret)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
