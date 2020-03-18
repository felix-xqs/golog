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
	xLogger Logger
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
		xLogger.h = NewHandlers(fileHandler)
	} else {
		fileHandler := NewFileHandler(c.FileConfig)
		xLogger.h = NewHandlers(fileHandler)
	}

}

func NewLogger(c *Config) *Logger {
	var xLogger = &Logger{}
	if c == nil {
		fileConfig := &FileConfig{
			Console: true,
		}
		c = &Config{FileConfig: fileConfig}
	}
	fileHandler := NewFileHandler(c.FileConfig)
	xLogger.SetHandlers(NewHandlers(fileHandler))
	return xLogger
}

func (xLogger *Logger) SetHandlers(hs *Handlers) {
	xLogger.h = hs
}

func DebugW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), DebugLevel, msg, m)
}

func InfoW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), InfoLevel, msg, m)
}

func WarnW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), WarnLevel, msg, m)
}

func ErrorW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), ErrorLevel, msg, m)
}

func FatalW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), FatalLevel, msg, m)
}

func PanicW(msg string, m map[string]interface{}) {
	xLogger.h.LogWith(context.Background(), PanicLevel, msg, m)
}

func Debug(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), DebugLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), InfoLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), ErrorLevel, format, args...)
}

func Panic(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), PanicLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	xLogger.h.Log(context.Background(), FatalLevel, format, args...)
}
