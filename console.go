package golog

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewConsoleFileHandler(c *FileConfig) *FileConsoleHandler {
	hook := &lumberjack.Logger{
		Filename:   c.LogFilePath + "info.log",
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
	}
	jsonInfo := zapcore.AddSync(hook)
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = ""
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = ""
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	xLogEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	var allCore []zapcore.Core
	infoCore := zapcore.NewCore(xLogEncoder, jsonInfo, infoPriority)
	allCore = append(allCore, infoCore)

	core := zapcore.NewTee(allCore...)
	var opts []zap.Option
	logger := zap.New(core).WithOptions(opts...)
	defer logger.Sync()
	return &FileConsoleHandler{logger: logger}
}

func consoleLog(logger *zap.Logger, msg string) {
	logger.Info(msg)
}
func (fch *FileConsoleHandler) Log(ctx context.Context, l Level, msg string, args ...interface{}) {
	logger := fch.logger
	logger = logger.With()
	consoleLog(logger, msg)
}
func (fch *FileConsoleHandler) LogWith(ctx context.Context, l Level, msg string, m map[string]interface{}) {
	logger := fch.logger
	consoleLog(logger, msg)
}
func (fch *FileConsoleHandler) Close() (err error) {
	return
}
