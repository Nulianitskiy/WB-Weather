package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "pkg/logger/app.log",
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

	return &ZapLogger{logger: zap.New(core)}
}

func (z *ZapLogger) Info(message string, args ...interface{}) {
	z.logger.Info(message, toZapFields(args...)...)
}

func (z *ZapLogger) Debug(message string, args ...interface{}) {
	z.logger.Debug(message, toZapFields(args...)...)
}

func (z *ZapLogger) Error(message string, args ...interface{}) {
	z.logger.Error(message, toZapFields(args...)...)
}

func (z *ZapLogger) Fatal(message string, args ...interface{}) {
	z.logger.Fatal(message, toZapFields(args...)...)
}

func toZapFields(args ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case zap.Field:
			fields = append(fields, v)
		case string:
			fields = append(fields, zap.String("msg", v))
		case int:
			fields = append(fields, zap.Int("num", v))
		default:
			fields = append(fields, zap.Any("unknown", v))
		}
	}
	return fields
}
