package option

import (
	"fmt"
	"testing"
	"time"
)

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Conf.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Conf.Timeout = timeout
	}
}

func MaxConn(maxConn int) Option {
	return func(s *Server) {
		s.Conf.Maxconns = maxConn
	}
}

func NewServer(addr string, port int, options ...Option) *Server {
	srv := &Server{
		Addr: addr,
		Port: port,
		Conf: &Config{
			Protocol: "tcp",
			Timeout:  30 * time.Second,
			Maxconns: 10,
		},
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

func TestOption(t *testing.T) {
	options := []Option{Protocol("udp"), Timeout(60 * time.Second), MaxConn(80)}
	srv := NewServer("127.0.0.1", 3306, options...)
	fmt.Printf("%+v\n", srv)
}
