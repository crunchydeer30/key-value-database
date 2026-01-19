package network

import (
	"encoding/binary"
	"fmt"
	"io"
)

func BuildPacket(payload []byte) []byte {
	packet := make([]byte, 4+len(payload))
	//nolint:gosec
	binary.BigEndian.PutUint32(packet[0:4], uint32(len(payload)))
	copy(packet[4:], payload)
	return packet
}

func ParsePacket(reader io.Reader) ([]byte, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(reader, lengthBuf); err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	payload := make([]byte, length)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return payload, nil
}
