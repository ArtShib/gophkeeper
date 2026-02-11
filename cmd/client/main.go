package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtShib/gophkeeper/internal/client/service/auth"
	"github.com/ArtShib/gophkeeper/internal/client/service/secret"
	"github.com/ArtShib/gophkeeper/internal/client/service/sync"
	"github.com/ArtShib/gophkeeper/internal/client/ui"
	//"github.com/ArtShib/gophkeeper/internal/client/ui/handlers"  // ✅ РАСКОММЕНТИРОВАТЬ!

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	model  *ui.AppModel
	width  int
	height int
}

func initialModel() *app {
	userSvc := &auth.Auth{}
	secretSvc := &secret.SecretSvc{}
	syncSvc := &sync.SyncService{}
	autoSyncSvc := &sync.AutoSyncService{}

	model := ui.NewAppModel(userSvc, secretSvc, syncSvc, autoSyncSvc)
	return &app{model: model}
}

func (a *app) Init() tea.Cmd {
	return nil
}

//	func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//		var cmd tea.Cmd
//
//		// ✅ Размер окна
//		if msg, ok := msg.(tea.WindowSizeMsg); ok {
//			a.width, a.height = msg.Width, msg.Height
//			return a, nil
//		}
//
//		// ✅ Клавиши → handlers (БЕЗ проверки Quit!)
//		if msg, ok := msg.(tea.KeyMsg); ok {
//			cmd = ui.InitHandlers(a.model)(msg.String())
//		}
//
//		return a, cmd // ✅ handlers сами возвращают tea.Quit при необходимости!
//	}
func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ✅ Размер окна
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		a.width, a.height = msg.Width, msg.Height
		return a, nil
	}

	// ✅ ВСТРОЕННАЯ НАВИГАЦИЯ - 100% РАБОТАЕТ!
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab", "Tab":
			// TAB по полям
			if a.model.ActiveTab == 0 {
				if a.model.FocusedField < 3 {
					a.model.FocusedField++
				}
			} else {
				if a.model.FocusedField < 4 {
					a.model.FocusedField++
				}
			}

		case "shift+tab", "Shift+Tab":
			if a.model.FocusedField > 0 {
				a.model.FocusedField--
			}

		case "left", "h":
			if a.model.ActiveTab > 0 {
				a.model.ActiveTab--
				a.model.FocusedField = 0
			}

		case "right", "l":
			if a.model.ActiveTab < 1 {
				a.model.ActiveTab++
				a.model.FocusedField = 0
			}

		case "f1", "F1":
			a.model.ShowPassword = !a.model.ShowPassword

		case "enter", "Enter":
			if a.model.FocusedField >= 3 {
				a.model.Loading = true
				a.model.CurrentScreen = "main"
			}

		case "q", "Q", "ctrl+c":
			return a, tea.Quit

		case "esc":
			a.model.CurrentScreen = "welcome"
		}

		return a, nil // ✅ ПЕРЕРИСОВКА!
	}

	return a, nil
}

func (a *app) View() string {
	return ui.View(a.model, a.width, a.height)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n👋 Выход...")
		os.Exit(0)
	}()

	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("❌ Ошибка UI: %v", err)
	}
}
