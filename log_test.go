package golog

import "testing"

func TestDebug(t *testing.T) {
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
