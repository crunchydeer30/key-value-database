package network

import (
	"bytes"
	"testing"
)

func TestTCPClient_Integration(t *testing.T) {
	handler := func(payload []byte) []byte {
		return []byte(payload)
	}

	server, err := NewTCPServer("0.0.0.0:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	go server.Serve()
	//nolint:errcheck
	defer server.Close()

	client, err := NewTCPClient(server.listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	//nolint:errcheck
	defer client.Close()

	message := []byte("hello")
	response, err := client.Send(message)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	if !bytes.Equal(response, message) {
		t.Fatalf("expected %s, got %s", message, response)
	}
}
