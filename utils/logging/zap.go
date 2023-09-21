package logging

import (
	"go.mrchanchal.com/zaphandler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
)

func NewZap(level string, dev bool) (*zap.Logger, error) {
	zapDev := dev
	var config zap.Config
	if zapDev {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Sampling = nil
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		parsedLevel = zapcore.DebugLevel
	}
	config.Level = zap.NewAtomicLevelAt(parsedLevel)
	return config.Build()
}

func NewSLogger(zlog *zap.Logger) *slog.Logger {
	return slog.New(zaphandler.New(zlog))
}
