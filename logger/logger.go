package logger

// Logger interface definition
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	// Add other methods as needed
}

var log Logger

// Init initializes the default logger
func Init() {
	log = NewDefaultLogger()
}

// SetLogger allows users to set their own logger
func SetLogger(logger Logger) {
	log = logger
}

// Get returns the current logger
func Get() Logger {
	return log
}
