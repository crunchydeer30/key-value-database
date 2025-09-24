package compute

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

type parserTestCase struct {
	name        string
	input       string
	wantCommand CommandName
	wantArgs    []string
	wantError   bool
	errType     error
}

func TestParser_Parse(t *testing.T) {
	logger := zap.NewNop()
	parser, _ := NewParser(logger)

	tests := []parserTestCase{
		{
			name:        "valid SET command",
			input:       "SET weather_2_pm cold_moscow_weather",
			wantCommand: SET,
			wantArgs:    []string{"weather_2_pm", "cold_moscow_weather"},
		},
		{
			name:        "valid GET command",
			input:       "GET key123",
			wantCommand: GET,
			wantArgs:    []string{"key123"},
		},
		{
			name:        "valid DEL command",
			input:       "DEL key_to_delete",
			wantCommand: DEL,
			wantArgs:    []string{"key_to_delete"},
		},
		{
			name:      "unknown command",
			input:     "FOO arg",
			wantError: true,
			errType:   ErrUnknownCommand,
		},
		{
			name:      "SET with missing argument",
			input:     "SET key_only",
			wantError: true,
			errType:   ErrInvalidNumberOfArgs,
		},
		{
			name:      "GET with too many arguments",
			input:     "GET key extra",
			wantError: true,
			errType:   ErrInvalidNumberOfArgs,
		},
		{
			name:      "empty input",
			input:     "   ",
			wantError: true,
			errType:   ErrInvalidQuery,
		},
		{
			name:        "extra spaces between words",
			input:       "SET   key   value",
			wantCommand: SET,
			wantArgs:    []string{"key", "value"},
		},
		{
			name:        "arguments with / and *",
			input:       "SET /path/to/file value*",
			wantCommand: SET,
			wantArgs:    []string{"/path/to/file", "value*"},
		},
		{
			name:      "empty argument",
			input:     "SET key ",
			wantError: true,
			errType:   ErrInvalidNumberOfArgs,
		},
		{
			name:      "command in lowercase",
			input:     "set key value",
			wantError: true,
			errType:   ErrUnknownCommand,
		},
		{
			name:      "too many arguments",
			input:     "SET key value extra",
			wantError: true,
			errType:   ErrInvalidNumberOfArgs,
		},
		{
			name:        "spaces at start and end",
			input:       "   GET key   ",
			wantCommand: GET,
			wantArgs:    []string{"key"},
		},
		{
			name:      "command without arguments",
			input:     "GET",
			wantError: true,
			errType:   ErrInvalidNumberOfArgs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := parser.Parse(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("invalid error, expected \"%v\", got \"%v\"", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if q.Command != tt.wantCommand {
				t.Errorf("expected command %s, got %v", tt.wantCommand, q.Command)
			}

			if len(q.Args) != len(tt.wantArgs) {
				t.Fatalf("expected args %v, got %v", tt.wantArgs, q.Args)
			}

			for i := range q.Args {
				if q.Args[i] != tt.wantArgs[i] {
					t.Errorf("arg: %d, expected %s, got %s", i, tt.wantArgs[i], q.Args[i])
				}
			}
		})
	}

}
