package bmff

import (
	"fmt"
	"runtime"
)

// DebugFlag
var (
	debugFlag = false
	log       Logger
)

// DebugLogger sets the debug logger
// TODO: rewrite
func DebugLogger(logger Logger) {
	debugFlag = true
	log = logger
}

// Logger is a logging Interface
type Logger interface {
	Debug(format string, args ...interface{})
}

// STDLogger is the Standard Logger
type STDLogger struct {
}

// Debug logs debug information
func (std STDLogger) Debug(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Printf("\n")

}

// Debug
// Info
// Warn
// Error
// Level

func traceBox(b Box, b2 box) {
	name := trace()
	log.Debug("%s\t %s\t Called from: %s", b2, b, name)
}

func traceBoxWithFlags(b Box, b2 box, f Flags) {
	name := trace()
	log.Debug("%s\t %s\t %s\t Called from: %s", b2, b, f, name)
}

func traceBoxWithMsg(b box, msg string) {
	name := trace()
	log.Debug("%s\t %s\t Called from: %s", b, msg, name)
}

func trace() string {
	pc, _, _, ok := runtime.Caller(3)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		return details.Name()
	}
	return ""
}
