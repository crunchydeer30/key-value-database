package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crunchydeer30/key-value-database/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

func NewLogger(cfg *config.LoggerConfig) (*zap.Logger, error) {
	var writer zapcore.WriteSyncer
	var encoder zapcore.Encoder

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	atomicLevel := zap.NewAtomicLevelAt(level)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder = zapcore.NewJSONEncoder(encoderConfig)

	switch cfg.Output {
	case "stdout":
		writer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writer = zapcore.AddSync(os.Stderr)
	default:
		dir := filepath.Dir(cfg.Output)
		//nolint:gosec
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}

		//nolint:gosec
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		writer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writer, atomicLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger, nil
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	default:
		return zap.InfoLevel, ErrInvalidLogLevel
	}
}
