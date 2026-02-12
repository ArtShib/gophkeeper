package ui

import (
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/models"
)

type SyncCompleteMsg struct{}
type SecretCreatedMsg struct{ Secret models.Secret }
type SecretDeletedMsg struct{ ID string }
type SecretUpdatedMsg struct{ Secret models.Secret }
type LoginSuccessMsg struct {
	User   *models.User
	Crypto *crypto.CryptoService
}
type RegisterSuccessMsg struct{}

type SecretItem struct {
	Secret models.Secret
}

func (i SecretItem) Title() string {
	return i.Secret.Metadata.Name
}

func (i SecretItem) Description() string {
	var desc string
	switch i.Secret.Type {
	case models.TypeCredentials:
		desc = "🔐 Пароль"
	case models.TypeText:
		desc = "📝 Заметка"
	case models.TypeCard:
		desc = "💳 Карта"
	case models.TypeBinary:
		desc = "📁 Файл"
	default:
		desc = string(i.Secret.Type)
	}

	statusIcon := ""
	switch i.Secret.Status {
	case models.StatusModified:
		statusIcon = " ✏️"
	case models.StatusNew:
		statusIcon = " 🆕"
	case models.StatusDeleted:
		statusIcon = " 🗑️"
	case models.StatusSynced:
		statusIcon = " ☁️"
	}

	return desc + statusIcon
}

func (i SecretItem) FilterValue() string {
	return i.Secret.Metadata.Name + " " + string(i.Secret.Metadata.Type) // + " " + strings.Join(i.Secret.Metadata.Extra, " ")
}

type UserWithServiceMsg struct {
	Username string
	Password string
}

type UserResultMsg struct {
	Success bool
	Error   string
	User    *models.User
	crypto  *crypto.CryptoService
}
