package secret

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/storage"
)

type SecretStore struct {
	store storage.GenericStorage[models.Secret]
}

func New(db *sql.DB, driver models.DriverType) *SecretStore {
	return &SecretStore{
		store: storage.GenericStorage[models.Secret]{
			DB:     db,
			Table:  "data",
			Driver: driver,
		},
	}
}

func (s *SecretStore) AddSecret(
	ctx context.Context,
	secretID string,
	userId int64,
	typeSecret string,
	data []byte,
	metadata string,
	createdAT int64) error {
	const op = "storage.postgres.AddSecret"
	query := `
			INSERT INTO gophkeeper.data (secret_id, owner_id, type, data, metadata, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.store.DB.ExecContext(ctx, query, secretID, userId, typeSecret, data, metadata, createdAT)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetUserSecrets список секретов
func (s *SecretStore) GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error) {
	const op = "storage.postgres.GetUserSecrets"

	query := `
		SELECT secret_id, type, data, metadata, created_at, updated_at
       	FROM data
       	WHERE owner_id = $1 AND is_deleted = FALSE`

	rows, err := s.store.DB.QueryContext(ctx, query, userID)

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
func (s *SecretStore) UpdateSecret(
	ctx context.Context,
	secretID string,
	userID int64,
	typeSecret string,
	data []byte,
	metadata string,
	updatedAT int64) error {
	const op = "storage.postgres.AddSecret"

	query := `
		  UPDATE data 
          SET data = $1, metadata = $2, updated_at = $3 
          WHERE secret_id = $4 AND owner_id = $5 AND type = $6`

	_, err := s.store.DB.ExecContext(ctx, query, data, metadata, updatedAT, secretID, userID, typeSecret)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// MarkDeleteSecret установка признака удаления
func (s *SecretStore) MarkDeleteSecret(
	ctx context.Context,
	secretID string,
	userID int64,
	typeSecret string,
	updatedAT int64) error {
	const op = "storage.postgres.AddSecret"
	query := `
			UPDATE gophkeeper.data 
			SET is_deleted = True, updated_at = $1 
			WHERE secret_id = $2 AND owner_id = $3 AND type = $4`

	_, err := s.store.DB.ExecContext(ctx, query, updatedAT, secretID, userID, typeSecret)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SecretStore) SetStatus(ctx context.Context, secretID string, status string) error {
	const op = "storage.postgres.SetStatus"

	query := `
		  UPDATE data 
          SET status = $1 
          WHERE secret_id = $2`

	_, err := s.store.DB.ExecContext(ctx, query, secretID, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
