package exit

import "os"

// типы с которыми приложение может завершать работу
const (
	ExitSuccess      = 0
	ExitStoreError   = 1
	ExitGRPCSrvError = 2
	ExitConfigError  = 3
	ExitOtherError   = 4
)

// Code тип для передчи информации о типе выхода из приложения
type Code int

// Exit метод для организации завершения приложения
func (c Code) Exit() {
	os.Exit(int(c))
}
