package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/service/secret"
	"github.com/ArtShib/gophkeeper/internal/client/service/sync"
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ===== Модель создания нового секрета =====
type NewSecretModel struct {
	secretType    models.SecretType
	nameInput     textinput.Model
	descInput     textinput.Model
	tagsInput     textinput.Model
	urlInput      textinput.Model
	usernameInput textinput.Model
	passwordInput textinput.Model
	notesInput    textinput.Model
	titleInput    textinput.Model
	contentInput  textinput.Model
	step          int // 0=тип, 1=метаданные, 2=данные
	errorMsg      string
	User          *models.User
	secretService *secret.SecretSvc
	ctx           context.Context
	syncService   *sync.SyncService
}

// Список доступных типов секретов для циклического перебора
var secretTypes = []models.SecretType{
	models.TypeCredentials,
	models.TypeText,
	models.TypeCard,
}

func InitialNewSecretModel(ctx context.Context, secretService *secret.SecretSvc, user *models.User, syncService *sync.SyncService) NewSecretModel {
	name := textinput.New()
	name.Placeholder = "Название секрета"
	name.Focus()
	name.CharLimit = 50
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Описание (необязательно)"
	desc.Width = 40
	desc.CharLimit = 100

	tags := textinput.New()
	tags.Placeholder = "Теги через запятую"
	tags.Width = 40
	tags.CharLimit = 100

	url := textinput.New()
	url.Placeholder = "https://example.com"
	url.Width = 40
	url.CharLimit = 200

	username := textinput.New()
	username.Placeholder = "Логин/Email"
	username.Width = 40
	username.CharLimit = 100

	password := textinput.New()
	password.Placeholder = "Пароль"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.Width = 40
	password.CharLimit = 100

	notes := textinput.New()
	notes.Placeholder = "Заметки (необязательно)"
	notes.Width = 40
	notes.CharLimit = 500

	title := textinput.New()
	title.Placeholder = "Заголовок заметки"
	title.Width = 40
	title.CharLimit = 100

	content := textinput.New()
	content.Placeholder = "Содержание заметки"
	content.Width = 40
	content.CharLimit = 1000

	return NewSecretModel{
		secretType:    models.TypeCredentials, // По умолчанию пароль
		nameInput:     name,
		descInput:     desc,
		tagsInput:     tags,
		urlInput:      url,
		usernameInput: username,
		passwordInput: password,
		notesInput:    notes,
		titleInput:    title,
		contentInput:  content,
		step:          0,
		ctx:           ctx,
		secretService: secretService,
		User:          user,
		syncService:   syncService,
	}
}

func (m NewSecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m NewSecretModel) View() string {
	var b strings.Builder

	b.WriteString(DetailTitleStyle.Render("➕ Создание нового секрета"))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString("Выберите тип секрета:\n\n")
		types := []struct {
			icon string
			name string
			desc string
			typ  models.SecretType
		}{
			{"🔐", "Пароль", "Логины и пароли от сайтов", models.TypeCredentials},
			{"📝", "Заметка", "Текстовые заметки", models.TypeText},
			{"💳", "Банковская карта", "Данные банковских карт", models.TypeCard},
		}

		for _, t := range types {
			style := DetailValueStyle
			if t.typ == m.secretType {
				style = style.Copy().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
				b.WriteString("👉 ")
			} else {
				b.WriteString("   ")
			}
			b.WriteString(style.Render(fmt.Sprintf("%s %s - %s", t.icon, t.name, t.desc)))
			b.WriteString("\n")
		}

	case 1:
		b.WriteString(DetailSectionStyle.Render("📋 Метаданные"))
		b.WriteString("\n\n")

		b.WriteString(LabelStyle.Render("Название:"))
		b.WriteString("\n")
		b.WriteString(m.nameInput.View())

		b.WriteString("\n\n")
		b.WriteString(LabelStyle.Render("Описание:"))
		b.WriteString("\n")
		b.WriteString(m.descInput.View())

		b.WriteString("\n\n")
		b.WriteString(LabelStyle.Render("Теги:"))
		b.WriteString("\n")
		b.WriteString(m.tagsInput.View())

	case 2:
		b.WriteString(DetailSectionStyle.Render("🔒 Данные секрета"))
		b.WriteString("\n\n")

		switch m.secretType {
		case models.TypeCredentials:
			b.WriteString(LabelStyle.Render("URL/Адрес:"))
			b.WriteString("\n")
			b.WriteString(m.urlInput.View())

			b.WriteString("\n\n")
			b.WriteString(LabelStyle.Render("Логин/Email:"))
			b.WriteString("\n")
			b.WriteString(m.usernameInput.View())

			b.WriteString("\n\n")
			b.WriteString(LabelStyle.Render("Пароль:"))
			b.WriteString("\n")
			b.WriteString(m.passwordInput.View())

			b.WriteString("\n\n")
			b.WriteString(LabelStyle.Render("Заметки:"))
			b.WriteString("\n")
			b.WriteString(m.notesInput.View())

		case models.TypeText:
			b.WriteString(LabelStyle.Render("Заголовок:"))
			b.WriteString("\n")
			b.WriteString(m.titleInput.View())

			b.WriteString("\n\n")
			b.WriteString(LabelStyle.Render("Содержание:"))
			b.WriteString("\n")
			b.WriteString(m.contentInput.View())

		case models.TypeCard:
			b.WriteString(InfoStyle.Render("💳 Создание банковской карты пока не реализовано"))
			b.WriteString("\n\n")
			b.WriteString("Нажмите Esc для возврата к выбору типа\n")
		}
	}

	if m.errorMsg != "" {
		b.WriteString("\n\n" + ErrorStyle.Render("❌ "+m.errorMsg))
	}

	b.WriteString("\n\n")
	switch m.step {
	case 0:
		b.WriteString(HelpStyle.Render("↑/↓ - выбор типа • Enter - далее • Esc - отмена"))
	case 1:
		b.WriteString(HelpStyle.Render("Tab - переключение полей • Enter - далее • Esc - назад"))
	case 2:
		b.WriteString(HelpStyle.Render("Tab - переключение полей • Enter - сохранить • Esc - назад"))
	}

	return b.String()
}

func (m NewSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.step == 0 {
				return InitialSecretsModel(m.ctx, m.User, m.secretService, m.syncService), nil
			} else {
				m.step--
				// При возврате на шаг 1, фокусируемся на поле названия
				if m.step == 1 {
					m.nameInput.Focus()
				}
				return m, nil
			}

		case "up", "k":
			if m.step == 0 {
				// Циклический перебор типов вверх
				for i, t := range secretTypes {
					if t == m.secretType {
						if i == 0 {
							m.secretType = secretTypes[len(secretTypes)-1]
						} else {
							m.secretType = secretTypes[i-1]
						}
						break
					}
				}
			}
			return m, nil

		case "down", "j":
			if m.step == 0 {
				// Циклический перебор типов вниз
				for i, t := range secretTypes {
					if t == m.secretType {
						if i == len(secretTypes)-1 {
							m.secretType = secretTypes[0]
						} else {
							m.secretType = secretTypes[i+1]
						}
						break
					}
				}
			}
			return m, nil

		case "tab", "shift+tab":
			if m.step == 1 || m.step == 2 {
				m = m.focusNextField()
			}
			return m, nil

		case "enter":
			if m.step == 0 {
				m.step = 1
				m.nameInput.Focus()
				return m, nil
			} else if m.step == 1 {
				if strings.TrimSpace(m.nameInput.Value()) == "" {
					m.errorMsg = "Название обязательно"
					return m, nil
				}
				m.step = 2
				m = m.focusFirstField()
				return m, nil
			} else {
				return m.createSecret()
			}

		case "ctrl+c":
			return InitialSecretsModel(m.ctx, m.User, m.secretService, m.syncService), tea.Quit

		default:
			// Обновляем активное поле ввода
			switch m.step {
			case 1:
				switch {
				case m.nameInput.Focused():
					m.nameInput, cmd = m.nameInput.Update(msg)
				case m.descInput.Focused():
					m.descInput, cmd = m.descInput.Update(msg)
				case m.tagsInput.Focused():
					m.tagsInput, cmd = m.tagsInput.Update(msg)
				}
			case 2:
				switch m.secretType {
				case models.TypeCredentials:
					switch {
					case m.urlInput.Focused():
						m.urlInput, cmd = m.urlInput.Update(msg)
					case m.usernameInput.Focused():
						m.usernameInput, cmd = m.usernameInput.Update(msg)
					case m.passwordInput.Focused():
						m.passwordInput, cmd = m.passwordInput.Update(msg)
					case m.notesInput.Focused():
						m.notesInput, cmd = m.notesInput.Update(msg)
					}
				case models.TypeText:
					switch {
					case m.titleInput.Focused():
						m.titleInput, cmd = m.titleInput.Update(msg)
					case m.contentInput.Focused():
						m.contentInput, cmd = m.contentInput.Update(msg)
					}
				}
			}
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// focusNextField - переключение на следующее поле ввода
func (m NewSecretModel) focusNextField() NewSecretModel {
	switch m.step {
	case 1:
		switch {
		case m.nameInput.Focused():
			m.nameInput.Blur()
			m.descInput.Focus()
		case m.descInput.Focused():
			m.descInput.Blur()
			m.tagsInput.Focus()
		case m.tagsInput.Focused():
			m.tagsInput.Blur()
			m.nameInput.Focus()
		}
	case 2:
		switch m.secretType {
		case models.TypeCredentials:
			switch {
			case m.urlInput.Focused():
				m.urlInput.Blur()
				m.usernameInput.Focus()
			case m.usernameInput.Focused():
				m.usernameInput.Blur()
				m.passwordInput.Focus()
			case m.passwordInput.Focused():
				m.passwordInput.Blur()
				m.notesInput.Focus()
			case m.notesInput.Focused():
				m.notesInput.Blur()
				m.urlInput.Focus()
			}
		case models.TypeText:
			switch {
			case m.titleInput.Focused():
				m.titleInput.Blur()
				m.contentInput.Focus()
			case m.contentInput.Focused():
				m.contentInput.Blur()
				m.titleInput.Focus()
			}
		}
	}
	return m
}

// focusFirstField - фокус на первом поле ввода текущего шага
func (m NewSecretModel) focusFirstField() NewSecretModel {
	switch m.step {
	case 1:
		m.nameInput.Focus()
		m.descInput.Blur()
		m.tagsInput.Blur()
	case 2:
		switch m.secretType {
		case models.TypeCredentials:
			m.urlInput.Focus()
			m.usernameInput.Blur()
			m.passwordInput.Blur()
			m.notesInput.Blur()
		case models.TypeText:
			m.titleInput.Focus()
			m.contentInput.Blur()
		}
	}
	return m
}

// createSecret - создание секрета из введенных данных
func (m NewSecretModel) createSecret() (tea.Model, tea.Cmd) {

	// Валидация обязательных полей
	if strings.TrimSpace(m.nameInput.Value()) == "" {
		m.errorMsg = "Название обязательно"
		return m, nil
	}

	// Создаем метаданные
	metadata := models.SecretMetadata{
		Name:  strings.TrimSpace(m.nameInput.Value()),
		Type:  models.SecretType(strings.TrimSpace(m.descInput.Value())),
		Extra: make(map[string]string),
	}

	// Добавляем теги в Extra, если они есть
	if tags := strings.TrimSpace(m.tagsInput.Value()); tags != "" {
		metadata.Extra["tags"] = tags
	}

	id, _ := crypto.GenerateUUID()

	// Создаем секрет
	newSecret := models.Secret{
		ID:        id,
		UserID:    m.User.ID,
		Type:      m.secretType,
		Metadata:  metadata,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: 0,
		Status:    models.StatusNew,
	}
	//

	// Сериализуем данные в зависимости от типа секрета
	switch m.secretType {
	case models.TypeCredentials:
		creds := models.Credentials{
			Username: strings.TrimSpace(m.usernameInput.Value()),
			Password: strings.TrimSpace(m.passwordInput.Value()),
			Notes:    strings.TrimSpace(m.notesInput.Value()),
		}
		data, err := json.Marshal(creds)
		if err != nil {
			m.errorMsg = "Ошибка сериализации данных"
			return m, nil
		}
		newSecret.Data = data

	case models.TypeText:
		note := models.Note{
			Title:   strings.TrimSpace(m.titleInput.Value()),
			Content: strings.TrimSpace(m.contentInput.Value()),
		}
		data, err := json.Marshal(note)
		if err != nil {
			m.errorMsg = "Ошибка сериализации данных"
			return m, nil
		}
		newSecret.Data = data

	case models.TypeCard:
		m.errorMsg = "Создание банковских карт пока не поддерживается"
		return m, nil
	}

	if err := m.secretService.AddSecret(m.ctx, &newSecret); err != nil {
		m.errorMsg = "Ошибка записи секрета в базу"
		return m, nil
	}
	// Возвращаемся к списку секретов с сообщением о создании
	return InitialSecretsModel(m.ctx, m.User, m.secretService, m.syncService), func() tea.Msg {
		return SecretCreatedMsg{Secret: newSecret}
	}
}
