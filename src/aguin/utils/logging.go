package utils

import (
	"log"
	"os"
	"fmt"
)

func GetLogger(prefix string) *Logger {
	return New(prefix)
}

func New(prefix string) *Logger {
	l := new(Logger)
	l.logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags)
	return l
}

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Critical(format string, args ...interface{}) {
	l.logger.Print(fmt.Sprintf(format, args))
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Print(fmt.Sprintf(format, args))
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}

func (l *Logger) Notice(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}