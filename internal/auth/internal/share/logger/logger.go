package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

// Logger represents a logger instance
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger instance
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, format string, v ...interface{}) {
	l.log(LevelDebug, ctx, format, v...)
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, format string, v ...interface{}) {
	l.log(LevelInfo, ctx, format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, format string, v ...interface{}) {
	l.log(LevelWarn, ctx, format, v...)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, format string, v ...interface{}) {
	l.log(LevelError, ctx, format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(ctx context.Context, format string, v ...interface{}) {
	l.log(LevelFatal, ctx, format, v...)
	os.Exit(1)
}

// log logs a message with the given level
func (l *Logger) log(level Level, ctx context.Context, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	// Get request ID from context if available
	requestID := "unknown"
	if id, ok := ctx.Value("request_id").(string); ok {
		requestID = id
	}

	// Format the message
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format(time.RFC3339)
	levelName := levelNames[level]

	// Log the message
	l.logger.Printf("[%s] [%s] [%s] %s", timestamp, levelName, requestID, msg)
}

// WithContext creates a new logger with the given context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l
}

// Default logger instance
var defaultLogger = New(LevelInfo)

// SetLevel sets the logging level for the default logger
func SetLevel(level Level) {
	defaultLogger.level = level
}

// Debug logs a debug message using the default logger
func Debug(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.Debug(ctx, format, v...)
}

// Info logs an info message using the default logger
func Info(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.Info(ctx, format, v...)
}

// Warn logs a warning message using the default logger
func Warn(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.Warn(ctx, format, v...)
}

// Error logs an error message using the default logger
func Error(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.Error(ctx, format, v...)
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.Fatal(ctx, format, v...)
}
