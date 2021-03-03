package bmff

import "fmt"

// DebugFlag
var (
	debugFlag = false
	log       Logger
)

func DebugLogger(logger Logger) {
	debugFlag = true
	log = logger
}

// Logger
type Logger interface {
	Debug(format string, args ...interface{})
}

type STDLogger struct{}

func (std STDLogger) Debug(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}

// Debug
// Info
// Warn
// Error
// Level
