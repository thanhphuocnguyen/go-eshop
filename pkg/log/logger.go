package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	Level zerolog.Level
}

func NewLogger(level *zerolog.Level) *Logger {
	if level != nil {
		return &Logger{Level: *level}
	}
	return &Logger{Level: zerolog.InfoLevel}
}
func (logger *Logger) print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msgf("%v", args...)
}
func (logger *Logger) Debug(args ...interface{}) {
	logger.print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.print(zerolog.FatalLevel, args...)
}
