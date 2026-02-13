package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/service/secret"
	"github.com/ArtShib/gophkeeper/internal/client/service/sync"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SecretsModel struct {
	secrets    models.ArraySecret
	secretList list.Model
	selected   *models.Secret
	viewMode   string
	filter     string
	loading    bool
	errorMsg   string
	statusBar  string
	username   string
	user       *models.User

	secretService *secret.SecretSvc
	syncService   *sync.SyncService
	ctx           context.Context
}

func InitialSecretsModel(ctx context.Context, user *models.User, secretService *secret.SecretSvc, syncService *sync.SyncService) SecretsModel {
	testSecrets, _ := secretService.GetUserSecrets(ctx, user.ID)

	//testSecrets := []models.Secret{
	//	{
	//		ID:     "1",
	//		UserID: 41,
	//		Type:   models.TypeCredentials,
	//		Data: func() []byte {
	//			data, _ := json.Marshal(models.Credentials{
	//				Username: "user@gmail.com",
	//				Password: "securepassword123",
	//				Notes:    "Рабочая почта",
	//			})
	//			return data
	//		}(),
	//		Metadata: models.SecretMetadata{
	//			Name: "Gmail",
	//			Type: "Рабочая почта",
	//			Extra: map[string]string{
	//				"email": "work",
	//			},
	//		},
	//		CreatedAt: time.Now().Unix(),
	//		UpdatedAt: 0,
	//		Status:    models.StatusNew,
	//	},
	//	{
	//		ID:     "2",
	//		UserID: 41,
	//		Type:   models.TypeCard,
	//		Data: func() []byte {
	//			data, _ := json.Marshal(models.Card{
	//				Number:   "4111111111111111",
	//				Holder:   "IVAN IVANOV",
	//				ExpiryMM: "12",
	//				ExpiryYY: "25",
	//				CVV:      "123",
	//				Bank:     "Sberbank",
	//			})
	//			return data
	//		}(),
	//		Metadata: models.SecretMetadata{
	//			Name:  "Основная карта",
	//			Type:  "Карта для онлайн-платежей",
	//			Extra: map[string]string{"bank": "main"},
	//		},
	//		CreatedAt: time.Now().Unix() - 172800,
	//		UpdatedAt: time.Now().Unix() - 3600,
	//		Status:    models.StatusModified,
	//	},
	//	{
	//		ID:     "3",
	//		UserID: 41,
	//		Type:   models.TypeText,
	//		Data: func() []byte {
	//			data, _ := json.Marshal(models.Note{
	//				Title:   "Коды доступа",
	//				Content: "WiFi: HomeNetwork123\nVPN: vpnpass456\nGit: githubtoken789",
	//			})
	//			return data
	//		}(),
	//		Metadata: models.SecretMetadata{
	//			Name:  "Коды доступа",
	//			Type:  "Важные коды для разных сервисов",
	//			Extra: map[string]string{"codes": "network"},
	//		},
	//		CreatedAt: time.Now().Unix(),
	//		UpdatedAt: time.Now().Unix(),
	//		Status:    models.StatusModified,
	//	},
	//}

	items := make([]list.Item, len(testSecrets))
	for i, secr := range testSecrets {
		items[i] = SecretItem{Secret: secr}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🔐 Мои секреты"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ListTitleStyle
	l.SetShowHelp(true)

	return SecretsModel{
		secrets:       testSecrets,
		secretList:    l,
		viewMode:      "list",
		statusBar:     "Готово",
		username:      user.Login,
		user:          user,
		secretService: secretService,
		syncService:   syncService,
		ctx:           ctx,
	}
}

func (m SecretsModel) Init() tea.Cmd {
	return nil
}

func (m SecretsModel) View() string {
	if m.loading {
		return "Загрузка секретов...\n"
	}

	var b strings.Builder

	switch m.viewMode {
	case "list":
		b.WriteString(TitleStyle.Render(fmt.Sprintf("🔐 GophKeeper - Привет, %s!", m.username)))
		b.WriteString("\n\n")

		b.WriteString(StatusBarStyle.Render(fmt.Sprintf(
			"Секретов: %s • Фильтр: %s • %s",
			SecretCountStyle.Render(fmt.Sprintf("%d", len(m.secrets))),
			m.filter,
			m.statusBar,
		)))
		b.WriteString("\n\n")

		b.WriteString(m.secretList.View())

		if m.errorMsg != "" {
			b.WriteString("\n" + ErrorStyle.Render(m.errorMsg))
		}

		b.WriteString("\n\n")
		b.WriteString("Enter - просмотр • N - новый секрет • D - удалить • F - фильтр • S - синхронизировать • L - выйти\n")

	case "detail":
		if m.selected == nil {
			m.viewMode = "list"
			return m.View()
		}

		b.WriteString(DetailTitleStyle.Render("🔍 Просмотр секрета"))
		b.WriteString("\n\n")

		b.WriteString(DetailSectionStyle.Render(
			DetailFieldStyle.Render("Название:") + m.selected.Metadata.Name + "\n" +
				DetailFieldStyle.Render("Тип:") + string(m.selected.Type) + "\n" +
				DetailFieldStyle.Render("Описание:") + string(m.selected.Metadata.Type) + "\n" +
				//DetailFieldStyle.Render("Теги:") + strings.Join(m.selected.Metadata.Extra, ", ") + "\n" +
				DetailFieldStyle.Render("Создан:") + time.Unix(m.selected.CreatedAt, 0).Format("02.01.2006 15:04") + "\n" +
				DetailFieldStyle.Render("Обновлен:") + time.Unix(m.selected.UpdatedAt, 0).Format("02.01.2006 15:04") + "\n" +
				DetailFieldStyle.Render("Статус:") + string(m.selected.Status),
		))

		b.WriteString("\n")
		b.WriteString(DetailTitleStyle.Render("📋 Данные секрета"))
		b.WriteString("\n\n")

		var dataStr string
		switch m.selected.Type {
		case models.TypeCredentials:
			var passData models.Credentials
			dataDecrypt, err := m.secretService.CryptoSvc.Decrypt(m.ctx, m.selected.Data)
			if err == nil {
				if err := json.Unmarshal(dataDecrypt, &passData); err == nil {
					dataStr = fmt.Sprintf(
						"%sСайт/Приложение:%s %s\n%sЛогин:%s %s\n%sПароль:%s %s\n%sЗаметки:%s %s",
						DetailFieldStyle.Render(), DetailValueStyle.Render(), passData.URL,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), passData.Username,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), "••••••••",
						DetailFieldStyle.Render(), DetailValueStyle.Render(), passData.Notes,
					)
				}
			}
		case models.TypeText:
			var noteData models.Note
			dataDecrypt, err := m.secretService.CryptoSvc.Decrypt(m.ctx, m.selected.Data)
			if err == nil {
				if err := json.Unmarshal(dataDecrypt, &noteData); err == nil {
					dataStr = fmt.Sprintf(
						"%sЗаголовок:%s %s\n%sСодержание:%s\n%s",
						DetailFieldStyle.Render(), DetailValueStyle.Render(), noteData.Title,
						DetailFieldStyle.Render(), DetailValueStyle.Render(),
						noteData.Content,
					)
				}
			}
		case models.TypeCard:
			var cardData models.Card
			dataDecrypt, err := m.secretService.CryptoSvc.Decrypt(m.ctx, m.selected.Data)
			if err == nil {
				if err := json.Unmarshal(dataDecrypt, &cardData); err == nil {
					maskedNumber := "•••• •••• •••• "
					if len(cardData.Number) >= 4 {
						maskedNumber += cardData.Number[len(cardData.Number)-4:]
					}
					dataStr = fmt.Sprintf(
						"%sБанк:%s %s\n%sНомер карты:%s %s\n%sДержатель:%s %s\n%sСрок действия:%s %s/%s\n%sCVV:%s %s",
						DetailFieldStyle.Render(), DetailValueStyle.Render(), cardData.Bank,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), maskedNumber,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), cardData.Holder,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), cardData.ExpiryMM, cardData.ExpiryYY,
						DetailFieldStyle.Render(), DetailValueStyle.Render(), "•••",
					)
				}
			}
		default:
			dataStr = DetailValueStyle.Render("Данные зашифрованы")
		}

		b.WriteString(dataStr)
		b.WriteString("\n\n")
		b.WriteString("Esc - назад • E - редактировать • C - копировать • D - удалить\n")
	}

	return b.String()
}

func (m SecretsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.viewMode {
		case "list":
			switch msg.String() {
			case "esc", "ctrl+c":
				return m, tea.Quit

			case "l", "L":
				return InitialChoiceModel(m.ctx, nil, nil, nil, nil, nil), nil

			case "enter":
				if selected, ok := m.secretList.SelectedItem().(SecretItem); ok {
					m.selected = &selected.Secret
					m.viewMode = "detail"
				}
				return m, nil

			case "n", "N":
				return InitialNewSecretModel(m.ctx, m.secretService, m.user, m.syncService), nil

			case "d", "D":
				if selected, ok := m.secretList.SelectedItem().(SecretItem); ok {
					m.statusBar = fmt.Sprintf("Удаление секрета '%s'...", selected.Secret.Metadata.Name)
					return m, func() tea.Msg {
						time.Sleep(1 * time.Second)
						return SecretDeletedMsg{ID: selected.Secret.ID}
					}
				}
				return m, nil

			case "s", "S":
				m.statusBar = "Синхронизация с сервером..."
				m.syncService.Sync(m.ctx)
				return InitialSecretsModel(m.ctx, m.user, m.secretService, m.syncService), nil

			case "f", "F":
				m.secretList.SetFilteringEnabled(!m.secretList.FilteringEnabled())
				if m.secretList.FilteringEnabled() {
					m.statusBar = "Режим фильтрации включен"
				} else {
					m.statusBar = "Режим фильтрации выключен"
				}
				return m, nil

			default:
				m.secretList, cmd = m.secretList.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "detail":
			switch msg.String() {
			case "esc":
				m.viewMode = "list"
				m.selected = nil
				return m, nil

			case "e", "E":
				if m.selected != nil {
					return InitialEditSecretModel(m.ctx, *m.selected, m.user, m.secretService, m.syncService), nil
				}
				return m, nil

			case "c", "C":
				m.statusBar = "Данные скопированы в буфер обмена"
				return m, nil

			case "d", "D":
				m.statusBar = fmt.Sprintf("Удаление секрета '%s'...", m.selected.Metadata.Name)
				return m, func() tea.Msg {
					time.Sleep(1 * time.Second)
					return SecretDeletedMsg{ID: m.selected.ID}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.secretList.SetSize(msg.Width, msg.Height-10)

	case SyncCompleteMsg:
		m.statusBar = "Синхронизация завершена"
		return m, nil

	case SecretDeletedMsg:
		for i, secret := range m.secrets {
			if secret.ID == msg.ID {
				m.secrets = append(m.secrets[:i], m.secrets[i+1:]...)
				break
			}
		}
		items := make([]list.Item, len(m.secrets))
		for i, secr := range m.secrets {
			items[i] = SecretItem{Secret: secr}
		}
		m.secretList.SetItems(items)
		m.statusBar = "Секрет удален"
		m.viewMode = "list"
		m.selected = nil
		return m, nil

	case SecretUpdatedMsg:
		for i, secr := range m.secrets {
			if secr.ID == msg.Secret.ID {
				m.secrets[i] = msg.Secret
				break
			}
		}
		items := make([]list.Item, len(m.secrets))
		for i, secr := range m.secrets {
			items[i] = SecretItem{Secret: secr}
		}
		m.secretList.SetItems(items)
		m.statusBar = "Секрет обновлен"
		m.viewMode = "list"
		m.selected = nil
		return m, nil

	case SecretCreatedMsg:
		m.secrets = append(m.secrets, msg.Secret)
		items := make([]list.Item, len(m.secrets))
		for i, secr := range m.secrets {
			items[i] = SecretItem{Secret: secr}
		}
		m.secretList.SetItems(items)
		m.statusBar = "Новый секрет создан"
		m.viewMode = "list"
		return m, nil

		//case tea.Model:
		//	if loginModel, ok := msg.(LoginModel); ok {
		//		return InitialSecretsModel(loginModel.username), nil
		//	}
	}

	return m, tea.Batch(cmds...)
}
