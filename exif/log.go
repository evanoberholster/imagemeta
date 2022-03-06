package exif

// Inital Log Level
var LogLevel = LogLevelError

const (
	LogLevelNone uint8 = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func checkLogLevel(l uint8) bool {
	return LogLevel <= l
}

func isError() bool {
	return checkLogLevel(LogLevelError)
}

func isWarn() bool {
	return checkLogLevel(LogLevelWarn)
}

func isInfo() bool {
	return checkLogLevel(LogLevelInfo)
}

func isDebug() bool {
	return checkLogLevel(LogLevelDebug)
}

func logTag() {

}

func logIfd() {

}
