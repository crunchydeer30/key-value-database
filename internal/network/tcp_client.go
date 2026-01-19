package network

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type TCPClient struct {
	address string
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func NewTCPClient(address string) (*TCPClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	return &TCPClient{
		address: address,
		conn:    conn,
		reader:  reader,
		writer:  writer,
	}, nil
}

func (c *TCPClient) Send(data []byte) ([]byte, error) {
	packet := make([]byte, 4+len(data))

	//nolint:gosec
	binary.BigEndian.PutUint32(packet[0:4], uint32(len(data)))

	copy(packet[4:], data)

	if _, err := c.writer.Write(packet); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	if err := c.writer.Flush(); err != nil {
		return nil, err
	}

	response, err := ParsePacket(c.reader)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("server closed connection")
		}

		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return response, nil
}

func (c *TCPClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
