package logger

// Logger is an interface for logging.
// It's just an example of how to use logger in the project.
// We can use any logger we want, for example, logrus, zap or zerolog,
// but I don't want to add any external dependencies.
type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Fatal(msg ...interface{})
}

const (
	// InfoLevel is a level for info messages.
	InfoLevel = "INFO"
	// ErrorLevel is a level for error messages.
	ErrorLevel = "ERROR"
	// DebugLevel is a level for debug messages.
	DebugLevel = "DEBUG"

	// fatalLevel is a level for fatalLevel messages.
	// It's used only for internal needs.
	fatalLevel = "FATAL"
)
