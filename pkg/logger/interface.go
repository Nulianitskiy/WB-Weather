package logger

type Logger interface {
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Error(massage string, args ...interface{})
	Fatal(massage string, args ...interface{})
}
