package network

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func encodeLength(length int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	return buf
}

func TestBuildPacket(t *testing.T) {
	tests := []struct {
		name        string
		description string
		data        []byte
	}{
		{
			name:        "valid small message",
			description: "should be built correctly",
			data: []byte{
				't', 'e', 's', 't',
			},
		},
		{
			name:        "empty payload",
			description: "should be built correctly",
			data:        []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := BuildPacket(tt.data)

			length := binary.BigEndian.Uint32(packet[0:4])
			if length != uint32(len(tt.data)) {
				t.Errorf("length header = %d, want %d", length, len(tt.data))
			}

			if string(packet[4:]) != string(tt.data) {
				t.Errorf("got %q, want %q", string(packet[4:]), string(tt.data))
			}
		})
	}
}

func TestParsePacket(t *testing.T) {
	tests := []struct {
		name        string
		description string
		length      int
		data        []byte
		wantError   bool
		wantLen     int
		wantStr     string
	}{
		{
			name:        "valid small message",
			description: "should be parsed correctly",
			length:      4,
			data: []byte{
				't', 'e', 's', 't',
			},
			wantError: false,
			wantLen:   4,
			wantStr:   "test",
		},
		{
			name:        "empty payload",
			description: "should be parsed correctly",
			length:      0,
			data:        []byte{},
			wantError:   false,
			wantLen:     0,
			wantStr:     "",
		},
		{
			name:        "payload longer than declared length",
			description: "should read only 4 bytes even though 5 are available",
			length:      4,
			data: []byte{
				't', 'e', 's', 't', '1',
			},
			wantError: false,
			wantLen:   4,
			wantStr:   "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append(encodeLength(tt.length), tt.data...)

			payload, err := ParsePacket(bytes.NewReader(data))

			if (err != nil) != tt.wantError {
				t.Errorf("error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(payload) != tt.wantLen {
					t.Errorf("len = %d, want %d", len(payload), tt.wantLen)
				}
				if string(payload) != tt.wantStr {
					t.Errorf("got %q, want %q", string(payload), tt.wantStr)
				}
			}
		})
	}
}
