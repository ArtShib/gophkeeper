package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type StoreSqlite struct {
	db *sql.DB
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

	return &StoreSqlite{db: db}, nil
}

func (s *StoreSqlite) Stop() error {
	return s.db.Close()
}

// AddUser регистрация пользователя
func (s *StoreSqlite) AddUser(ctx context.Context, id int64, login string, passHash []byte) error {
	const op = "storage.sqlite.AddUser"

	query := `
			INSERT INTO users (user_id, login, password_hash) 
			VALUES (?, ?, ?)
			ON CONFLICT(user_id) DO NOTHING`

	if _, err := s.db.ExecContext(ctx, query, id, login, passHash); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetUser select пользователя по логину
func (s *StoreSqlite) GetUser(ctx context.Context, login string, passHash []byte) (*models.User, error) {
	const op = "storage.sqlite.GetUser"

	query := `
			SELECT user_id, login, password_hash 
			FROM users 
			WHERE login = ? 
			  AND password_hash = ?`

	row := s.db.QueryRowContext(ctx, query, login, passHash)

	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// AddSecret добавление секрета
func (s *StoreSqlite) AddSecret(
	ctx context.Context,
	secretID string,
	userID int64,
	typeSecret string,
	data []byte,
	metadata string,
	createdAT int64,
	status string) error {
	const op = "storage.sqlite.AddSecret"
	query := `
			INSERT INTO secrets (secret_id, user_id, type, data, metadata, created_at, status) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, secretID, userID, typeSecret, data, metadata, createdAT, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetUserSecrets список секретов
func (s *StoreSqlite) GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error) {
	const op = "storage.sqlite.GetUserSecrets"

	query := `
			SELECT secret_id, type, data, metadata, created_at, updated_at, status
       		FROM secrets
       		WHERE user_id = ?
       			AND status != 'deleted'`

	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var secrets models.ArraySecret

	for rows.Next() {
		var secret models.Secret
		var updatedAt sql.NullInt64

		err := rows.Scan(&secret.ID, &secret.Type, &secret.Data, &secret.Metadata, &secret.CreatedAt, &updatedAt, &secret.Status)

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
func (s *StoreSqlite) UpdateSecret(
	ctx context.Context,
	secretID string,
	userID int64,
	typeSecret string,
	data []byte,
	metadata string,
	updatedAT int64,
	status string) error {
	const op = "storage.sqlite.AddSecret"

	query := `
			UPDATE secrets 
          	SET data = ?, metadata = ?, updated_at = ?, status = ? 
          	WHERE secret_id = ?
          	  AND user_id = ? 
          	  AND type = ?`

	_, err := s.db.ExecContext(ctx, query, data, metadata, updatedAT, status, secretID, userID, typeSecret)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteSecret установка признака удаления
func (s *StoreSqlite) deleteSecret(
	ctx context.Context,
	secretID string,
	ownerID int64) error {
	const op = "storage.sqlite.deleteSecret"
	query := `
			DELETE FROM secrets 
			WHERE secret_id = ? AND user_id = ?`

	_, err := s.db.ExecContext(ctx, query, secretID, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StoreSqlite) MarkDeleteSecret(
	ctx context.Context,
	secretID string,
	ownerID int64,
	updatedAT int64,
	status string) error {
	const op = "storage.sqlite.MarkDeleteSecret"
	query := `
			UPDATE secrets 
          	SET updated_at = ?, status = ? 
			WHERE secret_id = ? AND user_id = ?`

	_, err := s.db.ExecContext(ctx, query, updatedAT, status, secretID, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StoreSqlite) MarkSynced(ctx context.Context,
	secretID string,
	ownerID int64,
	isDeleted bool,
	status string) error {
	const op = "storage.sqlite.MarkSynced"
	if isDeleted {
		return s.deleteSecret(ctx, secretID, ownerID)
	}

	query := `
			UPDATE secrets 
          	SET status = ? 
			WHERE secret_id = ? AND user_id = ?`

	_, err := s.db.ExecContext(ctx, query, status, secretID, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StoreSqlite) GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error) {
	const op = "storage.sqlite.GetUserSecrets"

	query := `
			SELECT secret_id, type, data, metadata, created_at, updated_at, status
       		FROM secrets
       		WHERE user_id = ?
       			AND status != 'synced'`

	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var secrets models.ArraySecret

	for rows.Next() {
		var secret models.Secret
		var updatedAt sql.NullInt64

		err := rows.Scan(&secret.ID, &secret.Type, &secret.Data, &secret.Metadata, &secret.CreatedAt, &updatedAt, &secret.Status)

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
