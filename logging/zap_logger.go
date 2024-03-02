package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// createLoggerStd Метод инициализирует кастомный логгер с выводом в stdOut
func createLoggerStd() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

// createLoggerStdAndFile Метод инициализирует кастомный логгер с выводом в stdOut и в файл
// (lumberjack выполняет ротацию жерналов чтобы они не становились большими)
func createLoggerStdAndFile() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logging/info.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	logger := zap.New(core)

	// Добавляем базовое поле в логах
	return logger.With(
		zap.String("service", "Logger App"),
	)
}

func main() {
	logger := createLoggerStdAndFile()

	defer logger.Sync()

	logger.Info("Some Logging Message 1",
		zap.String("spanId", uuid.New().String()),
		zap.String("requestId", uuid.New().String()),
	)
	// {"level":"info","timestamp":"2024-03-02T19:47:01.221+0300","msg":"Some Logging Message 1","service":"Logger App","spanId":"cbbfac32-361f-4c4c-bd86-fe76c9aa1b1d","requestId":"a779cc8b-ad30-4979-a579-7bece402b081"}

	logger.Info("Some Logging Message 2")
	// {"level":"info","timestamp":"2024-03-02T19:47:01.232+0300","msg":"Some Logging Message 2","service":"Logger App"}
}
