package models

// SecretType алиас
type SecretType string
type SyncStatus string

const (
	TypeCredentials SecretType = "credentials"
	TypeText        SecretType = "text_data"
	TypeBinary      SecretType = "binary_data"
	TypeCard        SecretType = "credit_card"
	TypeUnspecifed  SecretType = "unspecifed"
	StatusNew       SyncStatus = "new"
	StatusModified  SyncStatus = "modified"
	StatusSynced    SyncStatus = "synced"
	StatusDeleted   SyncStatus = "deleted"
)

// Secret структура секрета
type Secret struct {
	ID        string         `json:"id"`
	UserID    int64          `json:"-"`
	Type      SecretType     `json:"type"`
	Data      []byte         `json:"data"`
	Metadata  SecretMetadata `json:"metadata"` //json.RawMessage SecretMetadata]
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	IsDeleted bool           `json:"is_deleted"`
	Status    SyncStatus     `json:"status"`
}

type SecretMetadata struct {
	Type  SecretType
	Name  string
	Extra map[string]string
}

// ArraySecret список секретов
type ArraySecret []Secret
