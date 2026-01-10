package logger

import (
	"errors"
	"testing"

	"github.com/crunchydeer30/key-value-database/internal/config"
)

func TestNewLogger(t *testing.T) {
	tests := map[string]struct {
		cfg           *config.Config
		expectedError error
	}{
		"valid config": {
			//nolint:exhaustruct
			cfg: &config.Config{
				Logger: config.LoggerConfig{
					Level:  "debug",
					Output: "stdout",
				},
			},
		},
		"invalid log level": {
			//nolint:exhaustruct
			cfg: &config.Config{
				Logger: config.LoggerConfig{
					Level:  "invalid",
					Output: "stdout",
				},
			},
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
