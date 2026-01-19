package config

import (
	"errors"
	"os"
	"testing"
)

func createTempConfigFile(t *testing.T, yml string) string {
	tmpFile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(yml); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

func createYmlConfig(engineType, loggerLevel, loggerOutput string) string {
	return "engine:\n" +
		"  type: " + engineType + "\n" +
		"logger:\n" +
		"  level: " + loggerLevel + "\n" +
		"  output: " + loggerOutput + "\n" +
		"network:\n"
}

func TestLoadValidConfig(t *testing.T) {
	engineType := "in_memory"
	loggerLevel := "debug"
	loggerOutput := "stdout"

	yml := createYmlConfig(engineType, loggerLevel, loggerOutput)

	path := createTempConfigFile(t, yml)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Engine.Type != engineType {
		t.Errorf("expected engine type %v, got %v", engineType, cfg.Engine.Type)
	}

	if cfg.Logger.Level != loggerLevel {
		t.Errorf("expected logger level %v, got %v", loggerLevel, cfg.Logger.Level)
	}
}

func TestLoadInvalidConfigWithValidationFailed(t *testing.T) {
	yml := createYmlConfig("in_memory", "invalid", "stdout")

	path := createTempConfigFile(t, yml)

	cfg, err := Load(path)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if cfg != nil {
		t.Errorf("expected nil, got %v", cfg)
	}

	if !errors.Is(err, ErrValidationFailed) {
		t.Errorf("expected error type %v, got %v", ErrValidationFailed, err)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	path := "invalid_path.yml"

	_, err := Load(path)
	if !errors.Is(err, ErrReadConfigFailed) {
		t.Errorf("expected error type %v, got %v", ErrReadConfigFailed, err)
	}
}
