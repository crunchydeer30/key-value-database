package inmemory

import (
	"errors"
	"testing"

	"github.com/crunchydeer30/key-value-database/internal/database/storage/engine"
	"go.uber.org/zap"
)

type getTestCase struct {
	name        string
	key         string
	wantValue   string
	wantError   bool
	wantErrType error
}

func newTestEngine(t *testing.T) engine.Engine {
	logger := zap.NewNop()
	e, err := NewInMemoryEngine(logger)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestInMemoryEngine_Get(t *testing.T) {
	tests := []getTestCase{
		{
			name:        "key does not exist",
			key:         "key2",
			wantValue:   "",
			wantError:   true,
			wantErrType: engine.ErrKeyNotFound,
		},
	}

	for _, tt := range tests {
		e := newTestEngine(t)
		t.Run(tt.name, func(t *testing.T) {
			value, err := e.Get(tt.key)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
				return
			}

			if value != tt.wantValue {
				t.Errorf("expected value %v, got %v", tt.wantValue, value)
			}
		})
	}

	t.Run("key exists", func(t *testing.T) {
		e := newTestEngine(t)
		if err := e.Set("key1", "value1"); err != nil {
			t.Fatal(err)
		}
		value, err := e.Get("key1")
		if err != nil {
			t.Fatal(err)
		}
		if value != "value1" {
			t.Errorf("expected value %v, got %v", "value1", value)
		}
	})
}

func TestInMemoryEngine_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "set new key",
			key:   "key1",
			value: "value1",
		},
		{
			name:  "set key with special chars",
			key:   "key_*/",
			value: "val123",
		},
		{
			name:  "overwrite existing key",
			key:   "key2",
			value: "value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEngine(t)

			if tt.name == "overwrite existing key" {
				if err := e.Set(tt.key, "old_value"); err != nil {
					t.Fatal(err)
				}
			}

			if err := e.Set(tt.key, tt.value); err != nil {
				t.Errorf("Set() error = %v", err)
			}

			got, err := e.Get(tt.key)
			if err != nil {
				t.Errorf("Get() after Set() error = %v", err)
			}
			if got != tt.value {
				t.Errorf("expected value %v, got %v", tt.value, got)
			}
		})
	}
}

func TestInMemoryEngine_Del(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setup     bool
		wantError bool
	}{
		{
			name:      "delete existing key",
			key:       "key1",
			setup:     true,
			wantError: false,
		},
		{
			name:      "delete non-existing key",
			key:       "key2",
			setup:     false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEngine(t)

			if tt.setup {
				if err := e.Set(tt.key, "value"); err != nil {
					t.Fatal(err)
				}
			}

			err := e.Del(tt.key)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got nil")
			}

			_, err = e.Get(tt.key)
			if tt.setup && !tt.wantError {
				if !errors.Is(err, engine.ErrKeyNotFound) {
					t.Errorf("expected key to be deleted, got error: %v", err)
				}
			}
		})
	}
}
