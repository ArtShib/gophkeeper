package sync

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/storage"
)

type SyncStore struct {
	store storage.GenericStorage[models.Secret]
}

func New(db *sql.DB, driver models.DriverType) *SyncStore {
	return &SyncStore{
		store: storage.GenericStorage[models.Secret]{
			DB:     db,
			Table:  "data",
			Driver: driver,
		},
	}
}

func (s *SyncStore) AddSecret(
	ctx context.Context,
	secretID string,
	userId int64,
	typeSecret string,
	data []byte,
	metadata string,
	createdAT int64) error {
	const op = "storage.AddSecret"
	query := `
			INSERT INTO data (secret_id, owner_id, type, data, metadata, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.store.DB.ExecContext(ctx, query, secretID, userId, typeSecret, data, metadata, createdAT)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *SyncStore) DeleteSecret(ctx context.Context, secretID string, userID int64) error {
	const op = "storage.deleteSecret"
	query := `
			DELETE FROM data 
			WHERE secret_id = $1 AND owner_id = $2`

	_, err := s.store.DB.ExecContext(ctx, query, secretID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *SyncStore) ListUserSecrets(ctx context.Context, userID int64) (models.ListSecrets, error) {
	const op = "storage.ListUserSecrets"

	countQuery := `SELECT COUNT(*) FROM data WHERE owner_id = $1`
	var count int
	err := s.store.DB.QueryRowContext(ctx, countQuery, userID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	query := `
			SELECT secret_id, updated_at, status
       		FROM data
       		WHERE owner_id = $1`

	rows, err := s.store.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	listSecrets := make(map[string]models.Secret, count)

	for rows.Next() {
		var secret models.Secret
		var updatedAt sql.NullInt64

		if err := rows.Scan(&secret.ID, &updatedAt, &secret.Status); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		if updatedAt.Valid {
			secret.UpdatedAt = updatedAt.Int64
		}
		listSecrets[secret.ID] = secret
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iter: %w", op, err)
	}

	return listSecrets, nil
}
func (s *SyncStore) GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error) {
	const op = "storage.GetSecretsToSync"

	query := `
			SELECT secret_id, type, data, metadata, created_at, updated_at, status
       		FROM data
       		WHERE owner_id = $1
       			AND status != 'synced'`

	rows, err := s.store.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var secrets models.ArraySecret

	for rows.Next() {
		var secret models.Secret
		var updatedAt sql.NullInt64
		var metadata string
		var secretType string
		err := rows.Scan(&secret.ID, &secretType, &secret.Data, &metadata, &secret.CreatedAt, &updatedAt, &secret.Status)

		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		if updatedAt.Valid {
			secret.UpdatedAt = updatedAt.Int64
		}
		if metadata != "" {
			if err := json.Unmarshal([]byte(metadata), &secret.Metadata); err != nil {
				return nil, fmt.Errorf("%s: unmarshal metadata: %w", op, err)
			}
		}
		secret.Type = models.SecretType(secretType)

		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iter: %w", op, err)
	}

	return secrets, nil
}
func (s *SyncStore) MarkSynced(ctx context.Context, secretID string, userID int64, status string) error {
	const op = "storage.MarkSynced"
	query := `
			UPDATE data 
          	SET status = $1 
			WHERE secret_id = $2 AND owner_id = $3`

	_, err := s.store.DB.ExecContext(ctx, query, status, secretID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *SyncStore) UpdateSecret(
	ctx context.Context,
	secretID string,
	userID int64,
	typeSecret string,
	data []byte,
	metadata string,
	updatedAT int64) error {
	const op = "storage.sync.UpdateSecret"

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
