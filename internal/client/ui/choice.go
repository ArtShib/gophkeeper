package ui

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ArtShib/gophkeeper/internal/client/config"
	"github.com/ArtShib/gophkeeper/internal/client/service/auth"

	"github.com/ArtShib/gophkeeper/internal/storage/secret"
	"github.com/ArtShib/gophkeeper/internal/storage/sync"
	tea "github.com/charmbracelet/bubbletea"
)

type ChoiceModel struct {
	cursor      int // 0 = Логин, 1 = Регистрация
	quitting    bool
	authService *auth.Auth
	ctx         context.Context
	config      *config.Config
	secretStore *secret.SecretStore
	syncStore   *sync.SyncStore
	logger      *slog.Logger
}

func InitialChoiceModel(ctx context.Context, authService *auth.Auth, config *config.Config, secretStore *secret.SecretStore, syncStore *sync.SyncStore, logger *slog.Logger) ChoiceModel {
	return ChoiceModel{
		cursor:      0,
		authService: authService,
		config:      config,
		secretStore: secretStore,
		syncStore:   syncStore,
		ctx:         ctx,
		logger:      logger,
	}
}

func (m ChoiceModel) Init() tea.Cmd {
	return nil
}

func (m ChoiceModel) View() string {
	if m.quitting {
		return "До свидания!\n"
	}

	var b strings.Builder

	b.WriteString(TitleStyle.Render("🔐 GophKeeper - Безопасное хранилище секретов"))
	b.WriteString("\n\n")

	b.WriteString("Добро пожаловать! Выберите действие:\n\n")

	if m.cursor == 0 {
		b.WriteString(SelectedMenuItemStyle.Render("> Вход в существующий аккаунт"))
	} else {
		b.WriteString(MenuItemStyle.Render("  Вход в существующий аккаунт"))
	}
	b.WriteString("\n")

	if m.cursor == 1 {
		b.WriteString(SelectedMenuItemStyle.Render("> Создать новый аккаунт"))
	} else {
		b.WriteString(MenuItemStyle.Render("  Создать новый аккаунт"))
	}

	b.WriteString("\n\n")
	b.WriteString("↑/↓ - выбор • Enter - подтвердить • Esc - выход\n")

	return b.String()
}

func (m ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			if m.cursor == 0 {
				return InitialLoginModel(m.ctx, m.authService, m.config, m.secretStore, m.syncStore, m.logger), nil
			} else {
				return InitialRegisterModel(m.ctx, m.authService, m.config, m.secretStore, m.syncStore, m.logger), nil
			}

		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}
