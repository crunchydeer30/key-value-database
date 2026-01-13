package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		return nil, err
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting read deadline: %v\n", err)
		return nil, err
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
	if _, err := c.writer.WriteString(string(data) + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
		return nil, err
	}

	if err := c.writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing writer: %v\n", err)
		return nil, err
	}

	response, err := c.reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		return nil, err
	}

	return response, nil
}

func (c *TCPClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
