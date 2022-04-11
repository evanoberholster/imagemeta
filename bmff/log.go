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

func traceBox(b Box, b2 box) {
	if log == nil {
		return
	}
	name := trace()
	log.Debug("%s\t %s\t Called from: %s", b2, b, name)
}

func tracebox(b box) {
	if log == nil {
		return
	}
	name := trace()
	log.Debug("%s\t Called from: %s", b, name)
}

func traceBoxWithFlags(b Box, b2 box, f Flags) {
	if log == nil {
		return
	}
	log.Debug("%s\t %s\t %s\t Called from: %s", b2, b, f, trace())
}

func traceBoxWithMsg(b box, msg string) {
	if log == nil {
		return
	}
	log.Debug("%s\t Called from: %s\n\t%s", b, trace(), msg)
}

func trace() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		return details.Name()
	}
	return ""
}
