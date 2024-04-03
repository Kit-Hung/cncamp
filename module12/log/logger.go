package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BaseLogger is a simplified abstraction of the zap.BaseLogger
type BaseLogger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) BaseLogger
}

// wrapper delegates all calls to the underlying zap.BaseLogger
type wrapper struct {
	logger *zap.Logger
}

// Debug logs an debug msg with fields
func (l wrapper) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l wrapper) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error msg with fields
func (l wrapper) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l wrapper) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l wrapper) With(fields ...zapcore.Field) BaseLogger {
	return wrapper{logger: l.logger.With(fields...)}
}
