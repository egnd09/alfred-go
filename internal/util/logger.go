package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

var (
	// Global logger instance
	std *Logger
)

// InitLogger initializes the global logger
func InitLogger(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	std = &Logger{logger.Sugar()}
	return std, nil
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if std == nil {
		// Fallback if not initialized
		logger, _ := zap.NewProduction()
		std = &Logger{logger.Sugar()}
	}
	return std
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

// With returns a logger with additional context
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(args...)}
}

// Named returns a logger with a specific name
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.SugaredLogger.Named(name)}
}

// NewStdoutLogger creates a development-friendly console logger
func NewStdoutLogger() (*Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger.Sugar()}, nil
}

// Fatal logs a message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalf(msg, args...)
	os.Exit(1)
}
