package network

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/crunchydeer30/key-value-database/internal/sync"
	"go.uber.org/zap"
)

type TCPServer struct {
	listener       net.Listener
	sem            *sync.Semaphore
	maxConnections int
	maxMessageSize uint32
	logger         *zap.Logger
	handler        Handler
}

type Handler func([]byte) []byte

func NewTCPServer(addr string, handler Handler, opts ...TCPServerOption) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", addr, err)
	}

	//nolint:exhaustruct
	s := &TCPServer{
		listener:       listener,
		handler:        handler,
		logger:         zap.NewNop(),
		maxMessageSize: 4096,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.maxConnections > 0 {
		s.sem = sync.NewSemaphore(s.maxConnections)
	}

	return s, nil
}

func (s *TCPServer) Serve() {
	s.logger.Info("started TCP server", zap.String("address", s.listener.Addr().String()))

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection", zap.Error(err))
			continue
		}

		if s.sem != nil {
			s.sem.Acquire()
		}

		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("recovered from panic", zap.Any("panic", r))
				}
			}()

			defer func() {
				if err := conn.Close(); err != nil {
					s.logger.Error("failed to close connection", zap.Error(err))
				}
			}()

			if s.sem != nil {
				defer s.sem.Release()
			}

			s.handle(conn)
		}(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	r := bufio.NewReader(conn)

	for {
		messageLengthBuffer := make([]byte, 4)
		if _, err := io.ReadFull(r, messageLengthBuffer); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			s.logger.Error("failed to read message length header", zap.Error(err))
			return
		}
		messageLength := binary.BigEndian.Uint32(messageLengthBuffer)

		if messageLength > uint32(s.maxMessageSize) {
			s.logger.Error("message too large",
				zap.Uint32("messageLength", messageLength),
				zap.Uint32("maxMessageSize", s.maxMessageSize),
			)
			return
		}

		payload := make([]byte, messageLength)
		if _, err := io.ReadFull(r, payload); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			s.logger.Error("failed to read message payload", zap.Error(err))
			return
		}

		result := s.handler(payload)

		responsePacket := make([]byte, 4+len(result))
		//nolint:gosec
		binary.BigEndian.PutUint32(responsePacket[0:4], uint32(len(result)))
		copy(responsePacket[4:], result)

		if _, err := conn.Write(responsePacket); err != nil {
			s.logger.Error("failed to write response", zap.Error(err))
			return
		}
	}
}
