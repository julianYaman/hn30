package utils

import (
	"log"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// Legacy logging functions - kept for backwards compatibility
// Consider migrating to structured logging (slog) for new code

func LogInfo(format string, v ...interface{}) {
	log.Printf(ColorGreen+"[INFO]"+ColorReset+" "+format, v...)
}

func LogWarn(format string, v ...interface{}) {
	log.Printf(ColorYellow+"[WARN]"+ColorReset+" "+format, v...)
}

func LogError(format string, v ...interface{}) {
	log.Printf(ColorRed+"[ERROR]"+ColorReset+" "+format, v...)
}

func LogComponent(component, format string, v ...interface{}) {
	log.Printf(ColorCyan+"[%s]"+ColorReset+" "+format, append([]interface{}{component}, v...)...)
}
