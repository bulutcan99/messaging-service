package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel = zapcore.DebugLevel + 2
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel + 2
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel = zapcore.WarnLevel + 2
	// ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel + 2
	// DPanicLevel logs are particularly important errors. In development the logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel + 2
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel + 2
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel + 2
)

func CreateLogger(selectedLevel int) *zap.Logger {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectRoot := filepath.Dir(filepath.Dir(rootPath))

	logPath := filepath.Join(projectRoot, "logs")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(err)
	}

	infoFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logPath, "info.log"),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     5, // days
		Compress:   true,
	})

	errorFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logPath, "error.log"),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     5, // days
		Compress:   true,
	})
	selected := InitLogger(selectedLevel)
	level := zap.NewAtomicLevelAt(selected)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	infoCore := zapcore.NewCore(fileEncoder, infoFile, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	}))

	errorCore := zapcore.NewCore(fileEncoder, errorFile, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	}))

	stdout := zapcore.AddSync(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		infoCore,
		errorCore,
	)

	return zap.New(core)
}

func InitLogger(level int) zapcore.Level {
	var loggerLevel zapcore.Level
	switch level {
	case int(DebugLevel):
		loggerLevel = zapcore.DebugLevel
	case int(InfoLevel):
		loggerLevel = zapcore.InfoLevel
	case int(WarnLevel):
		loggerLevel = zapcore.WarnLevel
	case int(ErrorLevel):
		loggerLevel = zapcore.ErrorLevel
	case int(DPanicLevel):
		loggerLevel = zapcore.DPanicLevel
	case int(PanicLevel):
		loggerLevel = zapcore.PanicLevel
	case int(FatalLevel):
		loggerLevel = zapcore.FatalLevel
	default:
		loggerLevel = zapcore.InfoLevel
	}
	return loggerLevel
}

func Infof(ctx context.Context, template string, args ...interface{}) {
	correlationID, ok := ctx.Value("correlation_id").(string)
	if !ok {
		zap.S().Infof(template, args...)
		return
	}

	zap.S().Infow(fmt.Sprintf(template, args...), zap.String("correlation_id", correlationID))
}

func Errorf(ctx context.Context, template string, args ...interface{}) {
	correlationID, ok := ctx.Value("correlation_id").(string)
	if !ok {
		zap.S().Errorf(template, args)
		return
	}

	zap.S().Errorw(fmt.Sprintf(template, args...), zap.String("correlation_id", correlationID))
}

func Warnf(ctx context.Context, template string, args ...interface{}) {
	correlationID, ok := ctx.Value("correlation_id").(string)
	if !ok {
		zap.S().Warnf(template, args...)
		return
	}

	zap.S().Warnw(fmt.Sprintf(template, args...), zap.String("correlation_id", correlationID))
}
