package logger

import (
	"fmt"
	"log"
)

type Logger struct {
	mod string
}

func For(name string) *Logger {
	return &Logger{mod: name}
}

func (l *Logger) WrapErr(label string, err error) error {
	if label == "" {
		return fmt.Errorf(l.mod+": %w", err)
	} else {
		return fmt.Errorf(l.mod+": "+label+": %w", err)
	}
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	log.Printf(l.mod+": "+msg, args...)
}

func (l *Logger) Println(msg string) {
	log.Println(l.mod + ": " + msg)
}
