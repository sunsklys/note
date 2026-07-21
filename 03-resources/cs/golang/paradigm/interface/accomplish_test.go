package _interface

import (
	"fmt"
	"testing"
)

type Shape interface {
	Sides() int
	Area() int
}
type Square struct {
	len int
}

func (s *Square) Sides() int {
	return 4
}

func (s *Square) Area() int {
	return s.len * s.len
}

func TestInterfaceCheck(t *testing.T) {
	s := Square{len: 5}
	fmt.Printf("%d\n", s.Sides())
	var _ Shape = (*Square)(nil)
}
