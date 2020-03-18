package golog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const TraceIDKey = "traceID"
const ServiceNameKey = "serviceName"

type FileHandler struct {
	logger      *zap.SugaredLogger
	level       Level
	serviceName string
}

type FileConfig struct {
	LogFilePath string `yaml:"logFilePath"` // 日志路径
	MaxSize     int    `yaml:"maxSize"`     // 单个日志最大的文件大小. 单位: MB
	MaxBackups  int    `yaml:"maxBackups"`  // 日志文件最多保存多少个备份
	MaxAge      int    `yaml:"maxAge"`      // 文件最多保存多少天
	Console     bool   `yaml:"console"`     // 是否命令行输出，开发环境可以使用
	LevelString string `yaml:"levelString"` // 输出的日志级别, 值：debug,info,warn,error,panic,fatal
	ServiceName string `yaml:"serviceName"` //服务名
}

func NewFileHandler(c *FileConfig) *FileHandler {
	xLogTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	hookInfo := lumberjack.Logger{
		Filename:   c.LogFilePath + "info.log",
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
	}

	hookError := lumberjack.Logger{
		Filename:   c.LogFilePath + "error.log",
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
	}

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	jsonErr := zapcore.AddSync(&hookError)
	jsonInfo := zapcore.AddSync(&hookInfo)

	// Optimize the xLog output for machine consumption and the console output
	// for human operators.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = ""
	encoderConfig.TimeKey = "time"
	//encoderConfig.CallerKey = "path" // 原定的path字段含义太多，建议还是分开，然后log调用的地方就叫caller
	encoderConfig.EncodeTime = xLogTimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	xLogEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var allCore []zapcore.Core
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the cores together.
	if c.LogFilePath != "" {
		errCore := zapcore.NewCore(xLogEncoder, jsonErr, errPriority)
		infoCore := zapcore.NewCore(xLogEncoder, jsonInfo, infoPriority)
		allCore = append(allCore, errCore, infoCore)
	}

	if c.Console {
		consoleDebugging := zapcore.Lock(os.Stdout)
		allCore = append(allCore, zapcore.NewCore(xLogEncoder, consoleDebugging, infoPriority))
	}

	core := zapcore.NewTee(allCore...)

	var opts []zap.Option

	logger := zap.New(core).WithOptions(opts...).Sugar()
	defer logger.Sync()

	return &FileHandler{logger: logger, level: LevelStringToCode(c.LevelString), serviceName: c.ServiceName}
}

func log(logger *zap.SugaredLogger, l Level, keysAndValues []interface{}) {
	switch l {
	case DebugLevel:
		logger.Debugw("", keysAndValues...)
	case InfoLevel:
		logger.Infow("", keysAndValues...)
	case WarnLevel:
		logger.Warnw("", keysAndValues...)
	case ErrorLevel:
		logger.Errorw("", keysAndValues...)
	case PanicLevel:
		logger.Panicw("", keysAndValues...)
	case FatalLevel:
		logger.Fatalw("", keysAndValues...)
	}
}

func (fh *FileHandler) Log(ctx context.Context, l Level, format string, args ...interface{}) {
	if l < fh.level {
		return
	}

	logger := fh.getLogger()

	traceID := ctx.Value(TraceIDKey)
	logger = logger.With(ServiceNameKey, fh.serviceName, TraceIDKey, traceID)

	msg := format
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	var keysAndValues []interface{}

	keysAndValues = append(keysAndValues, "msg")
	keysAndValues = append(keysAndValues, msg)

	log(logger, l, keysAndValues)
}
func NewConsoleFileHandler(c *FileConfig) *FileHandler {
	hook := &lumberjack.Logger{
		Filename:   c.LogFilePath,
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
	logger := zap.New(core).WithOptions(opts...).Sugar()
	defer logger.Sync()
	return &FileHandler{logger: logger, level: LevelStringToCode(c.LevelString), serviceName: c.ServiceName}
}
func (fh *FileHandler) LogWith(ctx context.Context, l Level, msg string, m map[string]interface{}) {
	if l < fh.level {
		return
	}

	logger := fh.getLogger()

	if msg != "" {
		m["msg"] = msg
	}
	if fh.serviceName != "" {
		m[ServiceNameKey] = fh.serviceName
	}

	var keysAndValues []interface{}

	for k, v := range m {
		keysAndValues = append(keysAndValues, k)
		keysAndValues = append(keysAndValues, v)
	}

	log(logger, l, keysAndValues)
}

func (fh *FileHandler) Close() (err error) {
	return
}

func getCaller(skip int) string {
	fileName, line, funcName := "???", 0, "???"
	pc, fileName, line, ok := runtime.Caller(skip)
	if ok {
		funcName = runtime.FuncForPC(pc).Name() // main.(*MyStruct).foo
		funcName = filepath.Base(funcName)      // .foo
		//funcName = strings.TrimPrefix(funcName, ".") // foo

		fileName = filepath.Base(fileName) // /full/path/basename.go => basename.go
	}

	ca := fileName + ":" + strconv.Itoa(line) + "(" + funcName + ")"
	//ca := fmt.Sprintf("%s:%d(%s)", fileName, line, funcName)
	return ca
}

func (fh *FileHandler) getLogger() (logger *zap.SugaredLogger) {
	logger = fh.logger.With("path", getCaller(5))
	return
}
