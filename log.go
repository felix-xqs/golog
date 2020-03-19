package golog

import (
	"context"
)

type Logger struct {
	h *Handlers
}

type Config struct {
	FileConfig *FileConfig `yaml:"fileConfig"`
}

var (
	logger Logger
)

func Init(c *Config) {
	if c == nil {
		fileConfig := &FileConfig{
			Console: true,
		}
		c = &Config{FileConfig: fileConfig}
	}
	if c.FileConfig.Console {
		fileHandler := NewConsoleFileHandler(c.FileConfig)
		logger.h = NewHandlers(fileHandler)
	} else {
		fileHandler := NewFileHandler(c.FileConfig)
		logger.h = NewHandlers(fileHandler)
	}

}

func NewLogger(c *Config) *Logger {
	var logger = &Logger{}
	if c == nil {
		fileConfig := &FileConfig{
			Console: true,
		}
		c = &Config{FileConfig: fileConfig}
	}
	if c.FileConfig.Console {
		fileHandler := NewConsoleFileHandler(c.FileConfig)
		logger.h = NewHandlers(fileHandler)
	} else {
		fileHandler := NewFileHandler(c.FileConfig)
		logger.h = NewHandlers(fileHandler)
	}
	return logger
}

func (logger *Logger) SetHandlers(hs *Handlers) {
	logger.h = hs
}

func Debug(format string, args ...interface{}) {
	logger.h.Log(context.Background(), DebugLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	logger.h.Log(context.Background(), InfoLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.h.Log(context.Background(), WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	logger.h.Log(context.Background(), ErrorLevel, format, args...)
}

func Panic(format string, args ...interface{}) {
	logger.h.Log(context.Background(), PanicLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.h.Log(context.Background(), FatalLevel, format, args...)
}

func DebugC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, DebugLevel, format, args...)
}

func InfoC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, InfoLevel, format, args...)
}

func WarnC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, WarnLevel, format, args...)
}

func ErrorC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, ErrorLevel, format, args...)
}

func PanicC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, PanicLevel, format, args...)
}

func FatalC(ctx context.Context, format string, args ...interface{}) {
	logger.h.Log(ctx, FatalLevel, format, args...)
}
