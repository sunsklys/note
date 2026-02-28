package option

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"
)

type ServerAll struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

func NewDefaultServer(addr string, port int) (*ServerAll, error) {
	return &ServerAll{addr, port, "tcp", 30 * time.Second, 100, nil}, nil
}

func NewTLSServer(addr string, port int, tls *tls.Config) (*ServerAll, error) {
	return &ServerAll{addr, port, "tcp", 30 * time.Second, 100, tls}, nil
}

func NewServerWithTimeout(addr string, port int, timeout time.Duration) (*ServerAll, error) {
	return &ServerAll{addr, port, "tcp", timeout, 100, nil}, nil
}

func NewTLSServerWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*ServerAll, error) {
	return &ServerAll{addr, port, "tcp", 30 * time.Second, maxconns, tls}, nil
}

func TestNew(t *testing.T) {
	fmt.Println(NewDefaultServer("127.0.0.1:8080", 8080))
	fmt.Println(NewTLSServer("127.0.0.1:8080", 8080, &tls.Config{}))
	fmt.Println(NewServerWithTimeout("127.0.0.1:8080", 8080, 30*time.Second))
	fmt.Println(NewTLSServerWithMaxConnAndTimeout("127.0.0.1:8080", 8080, 100, 30*time.Second, &tls.Config{}))
}
