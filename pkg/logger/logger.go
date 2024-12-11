package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		instance = initLogger()
	})

	return instance
}

func syncLogger() {
	if instance != nil {
		_ = instance.Sync()
	}
}

func initLogger() *zap.Logger {
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: "app.log",
		MaxSize:  10,
		MaxAge:   30,
		Compress: true,
	})

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, fileSyncer, zapcore.InfoLevel)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
