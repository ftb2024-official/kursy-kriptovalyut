package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func NewLogger() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename: "./../../logs/app.log",
			MaxSize:  10,
			MaxAge:   7,
		})
		level := zap.NewAtomicLevelAt(zap.InfoLevel)

		prodCfg := zap.NewProductionEncoderConfig()
		prodCfg.TimeKey = "ts"
		prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		devCfg := zap.NewDevelopmentEncoderConfig()
		devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(devCfg)
		fileEncoder := zapcore.NewJSONEncoder(prodCfg)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, level),
			zapcore.NewCore(fileEncoder, file, level),
		)

		logger = zap.New(core)
	})

	return logger
}
