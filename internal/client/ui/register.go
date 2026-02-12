package ui

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ArtShib/gophkeeper/internal/client/config"
	"github.com/ArtShib/gophkeeper/internal/client/service/auth"
	"github.com/ArtShib/gophkeeper/internal/storage/secret"
	"github.com/ArtShib/gophkeeper/internal/storage/sync"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type RegisterModel struct {
	usernameInput textinput.Model
	passwordInput textinput.Model
	confirmInput  textinput.Model
	errorMsg      string
	successMsg    string
	focus         int // 0=username, 1=password, 2=confirm
	submitting    bool
	authService   *auth.Auth
	config        *config.Config
	secretStore   *secret.SecretStore
	syncStore     *sync.SyncStore
	ctx           context.Context
	logger        *slog.Logger
}

func InitialRegisterModel(ctx context.Context, authService *auth.Auth, config *config.Config, secretStore *secret.SecretStore, syncStore *sync.SyncStore, logger *slog.Logger) RegisterModel {
	username := textinput.New()
	username.Placeholder = "Придумайте имя пользователя"
	username.Focus()
	username.CharLimit = 30
	username.Width = 40
	username.Prompt = "👤 "

	password := textinput.New()
	password.Placeholder = "Придумайте пароль"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.CharLimit = 30
	password.Width = 40
	password.Prompt = "🔒 "

	confirm := textinput.New()
	confirm.Placeholder = "Повторите пароль"
	confirm.EchoMode = textinput.EchoPassword
	confirm.EchoCharacter = '•'
	confirm.CharLimit = 30
	confirm.Width = 40
	confirm.Prompt = "✓ "

	return RegisterModel{
		usernameInput: username,
		passwordInput: password,
		confirmInput:  confirm,
		focus:         0,
		ctx:           ctx,
		authService:   authService,
		config:        config,
		secretStore:   secretStore,
		syncStore:     syncStore,
		logger:        logger,
	}
}

func (m RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("📝 Регистрация в GophKeeper"))
	b.WriteString("\n\n")

	b.WriteString(LabelStyle.Render("Имя пользователя:"))
	b.WriteString("\n")
	b.WriteString(m.usernameInput.View())

	b.WriteString("\n")
	b.WriteString(LabelStyle.Render("Пароль:"))
	b.WriteString("\n")
	b.WriteString(m.passwordInput.View())

	b.WriteString("\n")
	b.WriteString(LabelStyle.Render("Подтверждение пароля:"))
	b.WriteString("\n")
	b.WriteString(m.confirmInput.View())

	b.WriteString("\n\n")
	b.WriteString("Tab - переключение • Enter - регистрация • Esc - назад\n")

	if m.errorMsg != "" {
		b.WriteString("\n" + ErrorStyle.Render(m.errorMsg))
	}

	if m.successMsg != "" {
		b.WriteString("\n" + SuccessStyle.Render(m.successMsg))
	}

	if m.submitting {
		b.WriteString("\n⏳ Создание аккаунта...")
	}

	return b.String()
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			m.focus = (m.focus + 1) % 3
			m.usernameInput.Blur()
			m.passwordInput.Blur()
			m.confirmInput.Blur()

			switch m.focus {
			case 0:
				m.usernameInput.Focus()
			case 1:
				m.passwordInput.Focus()
			case 2:
				m.confirmInput.Focus()
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
			switch m.focus {
			case 0:
				m.usernameInput, cmd = m.usernameInput.Update(msg)
			case 1:
				m.passwordInput, cmd = m.passwordInput.Update(msg)
			case 2:
				m.confirmInput, cmd = m.confirmInput.Update(msg)
			}
			cmds = append(cmds, cmd)
		}

	case UserResultMsg:
		m.submitting = false

		if msg.Success && msg.User != nil {
			m.successMsg = "✅ Аккаунт успешно создан! Необходимо залогинется..."
			return InitialLoginModel(m.ctx, m.authService, m.config, m.secretStore, m.syncStore, m.logger), nil
		} else {
			m.errorMsg = msg.Error
			return m, nil
		}

	default:
		switch m.focus {
		case 0:
			m.usernameInput, cmd = m.usernameInput.Update(msg)
		case 1:
			m.passwordInput, cmd = m.passwordInput.Update(msg)
		case 2:
			m.confirmInput, cmd = m.confirmInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *RegisterModel) submit() (tea.Model, tea.Cmd) {
	username := strings.TrimSpace(m.usernameInput.Value())
	password := strings.TrimSpace(m.passwordInput.Value())
	confirm := strings.TrimSpace(m.confirmInput.Value())

	if username == "" || password == "" || confirm == "" {
		m.errorMsg = "Заполните все поля"
		return m, nil
	}

	if len(username) < 3 {
		m.errorMsg = "Имя пользователя должно быть не менее 3 символов"
		return m, nil
	}

	if len(password) < 6 {
		m.errorMsg = "Пароль должен быть не менее 6 символов"
		return m, nil
	}

	if password != confirm {
		m.errorMsg = "Пароли не совпадают"
		return m, nil
	}

	m.submitting = true
	m.errorMsg = ""
	m.successMsg = "⏳ Создание аккаунта..."

	return m, func() tea.Msg {
		user, err := m.authService.RegisterNewUser(m.ctx, username, password)
		if err != nil {
			return UserResultMsg{
				Success: false,
				Error:   err.Error(),
				User:    nil,
			}
		}

		return UserResultMsg{
			Success: true,
			Error:   "",
			User:    user,
		}
	}
}
