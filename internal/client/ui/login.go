package ui

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/config"
	mygrpc "github.com/ArtShib/gophkeeper/internal/client/grpc"
	"github.com/ArtShib/gophkeeper/internal/client/service/auth"
	svcSecret "github.com/ArtShib/gophkeeper/internal/client/service/secret"
	svcSync "github.com/ArtShib/gophkeeper/internal/client/service/sync"
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/storage/secret"
	"github.com/ArtShib/gophkeeper/internal/storage/sync"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	usernameInput textinput.Model
	passwordInput textinput.Model
	errorMsg      string
	successMsg    string
	focus         int // 0=username, 1=password
	submitting    bool
	username      string
	authService   *auth.Auth
	config        *config.Config
	secretStore   *secret.SecretStore
	syncStore     *sync.SyncStore
	ctx           context.Context
	logger        *slog.Logger
}

func InitialLoginModel(ctx context.Context, authService *auth.Auth, config *config.Config, secretStore *secret.SecretStore, syncStore *sync.SyncStore, logger *slog.Logger) LoginModel {
	username := textinput.New()
	username.Placeholder = "Введите имя пользователя"
	username.Focus()
	username.CharLimit = 30
	username.Width = 40
	username.Prompt = "👤 "

	password := textinput.New()
	password.Placeholder = "Введите пароль"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.CharLimit = 30
	password.Width = 40
	password.Prompt = "🔒 "

	return LoginModel{
		usernameInput: username,
		passwordInput: password,
		focus:         0,
		ctx:           ctx,
		authService:   authService,
		config:        config,
		secretStore:   secretStore,
		syncStore:     syncStore,
		logger:        logger,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("🔐 Вход в GophKeeper"))
	b.WriteString("\n\n")

	b.WriteString(LabelStyle.Render("Имя пользователя:"))
	b.WriteString("\n")
	b.WriteString(m.usernameInput.View())

	b.WriteString("\n")
	b.WriteString(LabelStyle.Render("Пароль:"))
	b.WriteString("\n")
	b.WriteString(m.passwordInput.View())

	b.WriteString("\n\n")
	b.WriteString("Tab - переключение • Enter - вход • Esc - назад\n")

	if m.errorMsg != "" {
		b.WriteString("\n" + ErrorStyle.Render(m.errorMsg))
	}

	if m.successMsg != "" {
		b.WriteString("\n" + SuccessStyle.Render(m.successMsg))
	}

	if m.submitting {
		b.WriteString("\n⏳ Авторизация...")
	}

	return b.String()
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			m.focus = (m.focus + 1) % 2
			if m.focus == 0 {
				m.usernameInput.Focus()
				m.passwordInput.Blur()
			} else {
				m.passwordInput.Focus()
				m.usernameInput.Blur()
			}
			return m, nil

		case "enter":
			if !m.submitting {
				return m.submit()
			}
			return m, nil

		case "esc":
			return InitialChoiceModel(m.ctx, m.authService, m.config, m.secretStore, m.syncStore, m.logger), nil

		case "ctrl+c":
			return m, tea.Quit

		default:
			if m.focus == 0 {
				m.usernameInput, cmd = m.usernameInput.Update(msg)
			} else {
				m.passwordInput, cmd = m.passwordInput.Update(msg)
			}
			cmds = append(cmds, cmd)
		}

	case UserResultMsg:
		m.submitting = false

		if msg.Success && msg.User != nil {
			m.successMsg = "✅ Успешный вход! Загрузка секретов..."
			m.username = msg.User.Login

			loginSuccess := LoginSuccessMsg{
				User:   msg.User,
				Crypto: msg.crypto,
			}
			return m, tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
				return loginSuccess
			})
		} else {
			m.errorMsg = msg.Error
			return m, nil
		}

	case LoginSuccessMsg:
		svc := svcSecret.New(m.secretStore, m.logger, msg.Crypto)
		syncSvcSecr := svcSync.NewSecretSvc(m.syncStore, m.logger, msg.Crypto)
		secretGRPC, _ := mygrpc.NewSecretClient(m.ctx, m.config.ConfigGRPC.Server, m.logger, m.config.ConfigTLS, msg.User.Token)
		syncSvc := svcSync.NewSyncServic(m.logger, syncSvcSecr, secretGRPC, msg.User.ID, msg.Crypto)
		return InitialSecretsModel(m.ctx, msg.User, svc, syncSvc), nil

	default:
		if m.focus == 0 {
			m.usernameInput, cmd = m.usernameInput.Update(msg)
		} else {
			m.passwordInput, cmd = m.passwordInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) submit() (tea.Model, tea.Cmd) {
	username := strings.TrimSpace(m.usernameInput.Value())
	password := strings.TrimSpace(m.passwordInput.Value())

	if username == "" || password == "" {
		m.errorMsg = "Заполните все поля"
		return m, nil
	}

	m.submitting = true
	m.errorMsg = ""
	m.successMsg = "⏳ Авторизация..."

	return m, func() tea.Msg {
		user, err := m.authService.Login(m.ctx, username, password)
		if err != nil {
			return UserResultMsg{
				Success: false,
				Error:   err.Error(),
				User:    nil,
				crypto:  nil,
			}
		}

		return UserResultMsg{
			Success: true,
			Error:   "",
			User:    user,
			crypto:  crypto.NewCryptoService(password, m.config.ConfigCrypto.SaltClient, m.logger),
		}
	}
}
