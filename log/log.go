package log

type Logger interface {
	Debug(logEntry *Entry)
	Info(logEntry *Entry)
	Warn(logEntry *Entry)
	Error(logEntry *Entry)
	Fatal(logEntry *Entry)
	Close() error
}
