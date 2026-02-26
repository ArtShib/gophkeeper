package ui

import "github.com/charmbracelet/lipgloss"

// ===== Стили =====
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			MarginBottom(1)

	SelectedMenuItemStyle = MenuItemStyle.Copy().
				Foreground(lipgloss.Color("#FF6B9D")).
				Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			MarginTop(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			MarginTop(1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45B7D1")).
			MarginBottom(1).
			MarginTop(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555")).
			Background(lipgloss.Color("#EEE")).
			Padding(0, 1).
			MarginTop(1)

	SecretCountStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	ListTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	DetailTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(lipgloss.Color("#7D56F4")).
				Bold(true).
				Padding(0, 2).
				MarginBottom(1)

	DetailFieldStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#45B7D1")).
				Bold(true).
				MarginRight(2)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#333")).
				MarginBottom(1)

	DetailSectionStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				MarginBottom(2)
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666")).
			Italic(true).
			MarginTop(1)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45B7D1")).
			Bold(true).
			MarginTop(1)
)
