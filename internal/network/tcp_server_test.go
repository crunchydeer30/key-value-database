package network

import (
	"bytes"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	readBuf  *bytes.Buffer
	writeBuf *bytes.Buffer
}

func newMockConn(incomingData []byte) *mockConn {
	return &mockConn{
		readBuf:  bytes.NewBuffer(incomingData),
		writeBuf: bytes.NewBuffer(nil),
	}
}

func (m *mockConn) Read(b []byte) (int, error)         { return m.readBuf.Read(b) }
func (m *mockConn) Write(b []byte) (int, error)        { return m.writeBuf.Write(b) }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestHandle_SingleMessage(t *testing.T) {
	handler := func(payload []byte) []byte {
		return []byte(payload)
	}

	server, err := NewTCPServer("0.0.0.0:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	incomingMessage := []byte("hello")
	incomingPacket := BuildPacket(incomingMessage)

	conn := newMockConn(incomingPacket)

	server.handle(conn)

	responsePacket := conn.writeBuf.Bytes()

	responseMessage, err := ParsePacket(bytes.NewReader(responsePacket))
	if err != nil {
		t.Fatalf("failed to parse response packet: %v", err)
	}

	if !bytes.Equal(incomingMessage, responseMessage) {
		t.Fatalf("expected %s, got %s", incomingMessage, responseMessage)
	}
}

func TestHandle_MultipleMessage(t *testing.T) {
	handler := func(payload []byte) []byte {
		return []byte(payload)
	}

	server, err := NewTCPServer("0.0.0.0:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	firstPacket := BuildPacket([]byte("first"))
	secondPacket := BuildPacket([]byte("second"))

	conn := newMockConn(append(firstPacket, secondPacket...))
	server.handle(conn)

	responsePackets := conn.writeBuf.Bytes()
	reader := bytes.NewReader(responsePackets)

	response1, err := ParsePacket(reader)
	if err != nil {
		t.Fatalf("failed to parse response packet: %v", err)
	}

	response2, err := ParsePacket(reader)
	if err != nil {
		t.Fatalf("failed to parse response packet: %v", err)
	}

	if !bytes.Equal(response1, []byte("first")) {
		t.Fatalf("expected %s, got %s", []byte("first"), response1)
	}

	if !bytes.Equal(response2, []byte("second")) {
		t.Fatalf("expected %s, got %s", []byte("second"), response2)
	}
}

func TestHandle_MalformedPacket(t *testing.T) {
	handler := func(payload []byte) []byte {
		return []byte(payload)
	}

	server, err := NewTCPServer("0.0.0.0:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	conn := newMockConn([]byte("malformed packet"))
	server.handle(conn)

	done := make(chan bool)

	go func() {
		server.handle(conn)
		done <- true
	}()

	select {
	case <-done:
		t.Log("handle() completed without hanging")
	case <-time.After(1 * time.Second):
		t.Fatal("handle() hung - timeout after 1 second")
	}

	if conn.writeBuf.Len() > 0 {
		t.Errorf("expected no response for malformed packet, got %d bytes", conn.writeBuf.Len())
	}
}
