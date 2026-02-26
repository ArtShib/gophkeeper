package ui

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/service/secret"
	"github.com/ArtShib/gophkeeper/internal/client/service/sync"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ===== Модель редактирования секрета =====
type EditSecretModel struct {
	secret          models.Secret
	nameInput       textinput.Model
	descInput       textinput.Model
	tagsInput       textinput.Model
	urlInput        textinput.Model
	usernameInput   textinput.Model
	passwordInput   textinput.Model
	notesInput      textinput.Model
	titleInput      textinput.Model
	contentInput    textinput.Model
	cardNumberInput textinput.Model
	cardHolderInput textinput.Model
	cardExpiryInput textinput.Model
	cardCVVInput    textinput.Model
	errorMsg        string
	step            int // 0=метаданные, 1=данные
	user            *models.User

	secretService *secret.SecretSvc
	ctx           context.Context
	syncService   *sync.SyncService
}

func InitialEditSecretModel(ctx context.Context, secret models.Secret, user *models.User, secretService *secret.SecretSvc, syncService *sync.SyncService) EditSecretModel {
	name := textinput.New()
	name.Placeholder = "Название секрета"
	name.SetValue(secret.Metadata.Name)
	name.Focus()
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Описание"
	desc.SetValue(secret.Metadata.Name)
	desc.Width = 40

	//tags := textinput.New()
	//tags.Placeholder = "Теги через запятую"
	//tags.SetValue(strings.Join(secret.Metadata.Extra, ", "))
	//tags.Width = 40

	url := textinput.New()
	url.Placeholder = "URL"
	url.Width = 40

	username := textinput.New()
	username.Placeholder = "Логин"
	username.Width = 40

	password := textinput.New()
	password.Placeholder = "Пароль"
	password.EchoMode = textinput.EchoPassword
	password.Width = 40

	notes := textinput.New()
	notes.Placeholder = "Заметки"
	notes.Width = 40

	title := textinput.New()
	title.Placeholder = "Заголовок"
	title.Width = 40

	content := textinput.New()
	content.Placeholder = "Содержание"
	content.Width = 40

	switch secret.Type {
	case models.TypeCredentials:
		var passData models.Credentials
		if err := json.Unmarshal(secret.Data, &passData); err == nil {
			url.SetValue(passData.URL)
			username.SetValue(passData.Username)
			password.SetValue(passData.Password)
			notes.SetValue(passData.Notes)
		}
	case models.TypeText:
		var noteData models.Note
		if err := json.Unmarshal(secret.Data, &noteData); err == nil {
			title.SetValue(noteData.Title)
			content.SetValue(noteData.Content)
		}
	}

	return EditSecretModel{
		secret:    secret,
		nameInput: name,
		descInput: desc,
		//tagsInput:     tags,
		urlInput:      url,
		usernameInput: username,
		passwordInput: password,
		notesInput:    notes,
		titleInput:    title,
		contentInput:  content,
		step:          0,
		user:          user,
		secretService: secretService,
		ctx:           ctx,
		syncService:   syncService,
	}
}

func (m EditSecretModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditSecretModel) View() string {
	var b strings.Builder

	b.WriteString(DetailTitleStyle.Render("✏️ Редактирование секрета"))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString(DetailSectionStyle.Render("Метаданные:\n\n"))
		b.WriteString("Название:\n")
		b.WriteString(m.nameInput.View())
		b.WriteString("\n\nОписание:\n")
		b.WriteString(m.descInput.View())
		b.WriteString("\n\nТеги:\n")
		b.WriteString(m.tagsInput.View())

	case 1:
		b.WriteString(DetailSectionStyle.Render("Данные секрета:\n\n"))
		switch m.secret.Type {
		case models.TypeCredentials:
			b.WriteString("URL:\n")
			b.WriteString(m.urlInput.View())
			b.WriteString("\n\nЛогин:\n")
			b.WriteString(m.usernameInput.View())
			b.WriteString("\n\nПароль:\n")
			b.WriteString(m.passwordInput.View())
			b.WriteString("\n\nЗаметки:\n")
			b.WriteString(m.notesInput.View())
		case models.TypeText:
			b.WriteString("Заголовок:\n")
			b.WriteString(m.titleInput.View())
			b.WriteString("\n\nСодержание:\n")
			b.WriteString(m.contentInput.View())
		case models.TypeCard:
			b.WriteString("💳 Редактирование банковской карты пока не реализовано\n")
		}
	}

	if m.errorMsg != "" {
		b.WriteString("\n" + ErrorStyle.Render(m.errorMsg))
	}

	b.WriteString("\n\n")
	if m.step == 0 {
		b.WriteString("Tab - переключение полей • Enter - далее • Esc - отмена\n")
	} else {
		b.WriteString("Tab - переключение полей • Enter - сохранить • Esc - назад\n")
	}

	return b.String()
}

func (m EditSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.step == 0 {
				return InitialSecretsModel(m.ctx, m.user, m.secretService, m.syncService), nil
			} else {
				m.step = 0
				m.nameInput.Focus()
				return m, nil
			}

		case "tab", "shift+tab":
			// Переключение между полями ввода
			m = m.focusNextField()
			return m, nil

		case "enter":
			if m.step == 0 {
				// Валидация на шаге метаданных
				if strings.TrimSpace(m.nameInput.Value()) == "" {
					m.errorMsg = "Название обязательно"
					return m, nil
				}
				// Переходим к шагу данных
				m.step = 1
				m.errorMsg = ""
				m = m.focusFirstField()
				return m, nil
			} else {
				// Сохраняем изменения
				return m.updateSecret()
			}

		case "ctrl+c":
			return InitialSecretsModel(m.ctx, m.user, m.secretService, m.syncService), tea.Quit

		default:
			// Обновляем активное поле ввода
			switch m.step {
			case 0:
				switch {
				case m.nameInput.Focused():
					m.nameInput, cmd = m.nameInput.Update(msg)
				case m.descInput.Focused():
					m.descInput, cmd = m.descInput.Update(msg)
				case m.tagsInput.Focused():
					m.tagsInput, cmd = m.tagsInput.Update(msg)
				}
				cmds = append(cmds, cmd)

			case 1:
				switch m.secret.Type {
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
					cmds = append(cmds, cmd)

				case models.TypeText:
					switch {
					case m.titleInput.Focused():
						m.titleInput, cmd = m.titleInput.Update(msg)
					case m.contentInput.Focused():
						m.contentInput, cmd = m.contentInput.Update(msg)
					}
					cmds = append(cmds, cmd)

				case models.TypeCard:
					switch {
					case m.cardNumberInput.Focused():
						m.cardNumberInput, cmd = m.cardNumberInput.Update(msg)
					case m.cardHolderInput.Focused():
						m.cardHolderInput, cmd = m.cardHolderInput.Update(msg)
					case m.cardExpiryInput.Focused():
						m.cardExpiryInput, cmd = m.cardExpiryInput.Update(msg)
					case m.cardCVVInput.Focused():
						m.cardCVVInput, cmd = m.cardCVVInput.Update(msg)
					}
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// focusNextField - переключение на следующее поле
func (m EditSecretModel) focusNextField() EditSecretModel {
	switch m.step {
	case 0:
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
		default:
			m.nameInput.Focus()
		}

	case 1:
		switch m.secret.Type {
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
			default:
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
			default:
				m.titleInput.Focus()
			}

		case models.TypeCard:
			switch {
			case m.cardNumberInput.Focused():
				m.cardNumberInput.Blur()
				m.cardHolderInput.Focus()
			case m.cardHolderInput.Focused():
				m.cardHolderInput.Blur()
				m.cardExpiryInput.Focus()
			case m.cardExpiryInput.Focused():
				m.cardExpiryInput.Blur()
				m.cardCVVInput.Focus()
			case m.cardCVVInput.Focused():
				m.cardCVVInput.Blur()
				m.cardNumberInput.Focus()
			default:
				m.cardNumberInput.Focus()
			}
		}
	}

	return m
}

// focusFirstField - фокус на первом поле текущего шага
func (m EditSecretModel) focusFirstField() EditSecretModel {
	switch m.step {
	case 0:
		m.nameInput.Focus()
		m.descInput.Blur()
		m.tagsInput.Blur()
	case 1:
		switch m.secret.Type {
		case models.TypeCredentials:
			m.urlInput.Focus()
			m.usernameInput.Blur()
			m.passwordInput.Blur()
			m.notesInput.Blur()
		case models.TypeText:
			m.titleInput.Focus()
			m.contentInput.Blur()
		case models.TypeCard:
			m.cardNumberInput.Focus()
			m.cardHolderInput.Blur()
			m.cardExpiryInput.Blur()
			m.cardCVVInput.Blur()
		}
	}
	return m
}

// updateSecret - обновление секрета в БД
func (m EditSecretModel) updateSecret() (tea.Model, tea.Cmd) {
	// Валидация
	if strings.TrimSpace(m.nameInput.Value()) == "" {
		m.errorMsg = "Название обязательно"
		return m, nil
	}

	// Обновляем метаданные
	m.secret.Metadata.Name = strings.TrimSpace(m.nameInput.Value())

	// Очищаем Extra и добавляем новые значения
	//m.secret.Metadata.Extra = make([]string, 0)
	//
	//// Добавляем описание
	//if desc := strings.TrimSpace(m.descInput.Value()); desc != "" {
	//	m.secret.Metadata.Extra = append(m.secret.Metadata.Extra, "description:"+desc)
	//}
	////
	////// Добавляем теги
	////if tags := strings.TrimSpace(m.tagsInput.Value()); tags != "" {
	////	tagList := strings.Split(tags, ",")
	////	cleanedTags := make([]string, 0)
	////	for _, tag := range tagList {
	////		if cleanedTag := strings.TrimSpace(tag); cleanedTag != "" {
	////			cleanedTags = append(cleanedTags, cleanedTag)
	////		}
	////	}
	////	if len(cleanedTags) > 0 {
	////		m.secret.Metadata.Extra = append(m.secret.Metadata.Extra, "tags:"+strings.Join(cleanedTags, ","))
	////	}
	////}

	// Обновляем данные в зависимости от типа секрета
	switch m.secret.Type {
	case models.TypeCredentials:
		creds := models.Credentials{
			URL:      strings.TrimSpace(m.urlInput.Value()),
			Username: strings.TrimSpace(m.usernameInput.Value()),
			Password: strings.TrimSpace(m.passwordInput.Value()),
			Notes:    strings.TrimSpace(m.notesInput.Value()),
		}
		data, err := json.Marshal(creds)
		if err != nil {
			m.errorMsg = "Ошибка сериализации данных: " + err.Error()
			return m, nil
		}
		m.secret.Data = data

	case models.TypeText:
		note := models.Note{
			Title:   strings.TrimSpace(m.titleInput.Value()),
			Content: strings.TrimSpace(m.contentInput.Value()),
		}
		data, err := json.Marshal(note)
		if err != nil {
			m.errorMsg = "Ошибка сериализации данных: " + err.Error()
			return m, nil
		}
		m.secret.Data = data

	case models.TypeCard:
		card := models.Card{
			Number:   strings.TrimSpace(m.cardNumberInput.Value()),
			Holder:   strings.TrimSpace(m.cardHolderInput.Value()),
			ExpiryYY: strings.TrimSpace(m.cardExpiryInput.Value()),
			CVV:      strings.TrimSpace(m.cardCVVInput.Value()),
		}
		data, err := json.Marshal(card)
		if err != nil {
			m.errorMsg = "Ошибка сериализации данных: " + err.Error()
			return m, nil
		}
		m.secret.Data = data
	}

	// Обновляем временные метки и статус
	m.secret.UpdatedAt = time.Now().Unix()
	m.secret.Status = models.StatusModified

	// Сохраняем изменения в БД
	if err := m.secretService.UpdateSecret(m.ctx, &m.secret); err != nil {
		m.errorMsg = "Ошибка обновления секрета: " + err.Error()
		return m, nil
	}

	// Запускаем синхронизацию в фоне
	if m.syncService != nil {
		go m.syncService.Sync(m.ctx)
	}

	// Возвращаемся к списку секретов с сообщением об успешном обновлении
	return InitialSecretsModel(m.ctx, m.user, m.secretService, m.syncService), func() tea.Msg {
		return SecretUpdatedMsg{Secret: m.secret}
	}
}
