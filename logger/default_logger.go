package logger

import "go.uber.org/zap"

type DefaultLogger struct {
	zapLogger *zap.Logger
}

// Error implements Logger.
func (l *DefaultLogger) Error(args ...interface{}) {
	panic("unimplemented")
}

// Warn implements Logger.
func (l *DefaultLogger) Warn(args ...interface{}) {
	panic("unimplemented")
}

func NewDefaultLogger() *DefaultLogger {
	logger, _ := zap.NewProduction()
	return &DefaultLogger{zapLogger: logger}
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	l.zapLogger.Sugar().Debug(args...)
}

func (l *DefaultLogger) Info(args ...interface{}) {
	l.zapLogger.Sugar().Info(args...)
}

// Implement other methods...
