package delete_models

//var (
//	ErrUserExists         = errors.New("user already exists")
//	ErrUserNotFound       = errors.New("user not found")
//	ErrInvalidCredentials = errors.New("invalid credentials")
//	ErrLoginOrPassIsEmpty = errors.New("login and password_hash required")
//)

// User структура
//type User struct {
//	ID           int64  `json:"id"`
//	Login        string `json:"login"`
//	PasswordHash []byte `json:"password_hash"`
//	Token        string
//}
//
//// SecretType алиас
//type SecretType string
//type SyncStatus string
//
//const (
//	TypeCredentials SecretType = "credentials"
//	TypeText        SecretType = "text_data"
//	TypeBinary      SecretType = "binary_data"
//	TypeCard        SecretType = "credit_card"
//	StatusNew       SyncStatus = "new"
//	StatusModified  SyncStatus = "modified"
//	StatusSynced    SyncStatus = "synced"
//	StatusDeleted   SyncStatus = "deleted"
//)
//
//// Secret структура секрета
//type Secret struct {
//	ID     string `json:"id"`
//	UserID int64  `json:"-"`
//	Type      SecretType     `json:"type"`
//	Data      []byte         `json:"data"`
//	Metadata  SecretMetadata `json:"metadata"`
//	CreatedAt int64          `json:"created_at"`
//	UpdatedAt int64          `json:"updated_at"`
//	IsDeleted bool           `json:"is_deleted"`
//	Status    SyncStatus     `json:"status"`
//}

//type SecretMetadata struct {
//	Type  SecretType
//	Name  string
//	Extra map[string]string
//}

//// ArraySecret список секретов
//type ArraySecret []Secret
//
//type contextKey string
//
//// contextKey
//const (
//	UserIDKey contextKey = "userID"
//)
