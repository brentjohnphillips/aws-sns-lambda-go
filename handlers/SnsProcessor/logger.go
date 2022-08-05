package main

import (
	"fmt"
)

type Level int64

func New() *Logger {
	return &Logger{
		Level: InfoLevel,
	}
}

const (
	ErrorLevel Level = iota
	WarnLevel
	DebugLevel
	InfoLevel
)

type Logger struct {
	Level Level
}

func (logger *Logger) SetLevel(level Level) {
	logger.Level = level
}
func (logger *Logger) IsLevelEnabled(level Level) bool {
	//Error is 0, Info = 3, so if the set log level is less than or equal to message level we can log

	return logger.Level >= level
}

func (logger *Logger) Log(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.Log(ErrorLevel, format, args...)
}
func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.Log(WarnLevel, format, args...)
}
func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.Log(DebugLevel, format, args...)
}
func (logger *Logger) Info(format string, args ...interface{}) {
	logger.Log(InfoLevel, format, args...)
}
