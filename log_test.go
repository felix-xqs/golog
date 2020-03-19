package golog

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestConsole(t *testing.T) {
	Init(&Config{FileConfig: &FileConfig{
		LogFilePath: "./log/1",
		MaxSize:     1,
		MaxBackups:  1,
		MaxAge:      1,
		Console:     true,
		LevelString: "info",
		ServiceName: "test",
	}})
	Info("string data")
}

func TestDebug(t *testing.T) {
	Init(&Config{FileConfig: &FileConfig{
		LogFilePath: "./log/2",
		MaxSize:     1,
		MaxBackups:  1,
		MaxAge:      1,
		LevelString: "debug",
		ServiceName: "test",
	}})
	Info("debug data %s", "debug")
}

func TestDebugC(t *testing.T) {
	Init(&Config{FileConfig: &FileConfig{
		LogFilePath: "./log/2",
		MaxSize:     1,
		MaxBackups:  1,
		MaxAge:      1,
		LevelString: "debug",
		ServiceName: "test",
	}})
	ctx := &gin.Context{
		Request:  nil,
		Writer:   nil,
		Params:   nil,
		Keys:     nil,
		Errors:   nil,
		Accepted: nil,
	}
	traceID := uuid.NewV4().String()
	ctx.Set(TraceIDKey, traceID)
	InfoC(ctx, "debug data %s", "debug")
}
