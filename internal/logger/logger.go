package logger

import (
	"log"
	"os"
	"strings"
)

type logger struct {
	log   *log.Logger
	level int
}

func NewLogger(level string) Logger {
	return &logger{
		log:   log.New(log.Writer(), "", log.Flags()),
		level: castLevel(level),
	}
}

func (l *logger) Debug(msg ...interface{}) {
	if l.level <= 1 {
		l.print(DebugLevel, msg...)
	}
}

func (l *logger) Info(msg ...interface{}) {
	if l.level <= 2 {
		l.print(InfoLevel, msg...)
	}
}

func (l *logger) Error(msg ...interface{}) {
	if l.level <= 3 {
		l.print(ErrorLevel, msg...)
	}
}

func (l *logger) Fatal(msg ...interface{}) {
	l.print(fatalLevel, msg...)
	os.Exit(1)
}

func (l *logger) print(level string, msg ...interface{}) {
	l.log.Println(append([]interface{}{level}, msg...)...)
}

// Cast level to string.
func castLevel(level string) int {
	switch strings.ToUpper(level) {
	case DebugLevel:
		return 1
	case InfoLevel:
		return 2
	case ErrorLevel:
		return 3
	default:
		return 2
	}
}
