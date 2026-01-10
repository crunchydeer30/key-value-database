package logger

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/crunchydeer30/key-value-database/internal/config"
)

func TestNewLogger(t *testing.T) {
	tests := map[string]struct {
		cfg           *config.LoggerConfig
		expectedError error
	}{
		"valid config": {
			//nolint:exhaustruct
			cfg: &config.LoggerConfig{
				Level:  "debug",
				Output: "stdout",
			},
		},
		"invalid log level": {
			//nolint:exhaustruct
			cfg: &config.LoggerConfig{
				Level:  "invalid",
				Output: "stdout",
			},
			expectedError: ErrInvalidLogLevel,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLogger(tt.cfg)
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if logger == nil {
				t.Errorf("expected logger, got nil")
			}
		})
	}
}

func TestNewLoggerCreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "subdir", "app.log")

	cfg := &config.LoggerConfig{
		Level:  "info",
		Output: logPath,
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}

	dirPath := filepath.Dir(logPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("directory %s was not created", dirPath)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("log file %s was not created", logPath)
	}

	logger.Info("test message")
	err = logger.Sync()
	if err != nil {
		t.Fatalf("failed to sync logger: %v", err)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "test message") {
		t.Error("log message was not written to file")
	}
}
