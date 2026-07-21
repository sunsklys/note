package err

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

var b = []byte{0x48, 0x61, 0x6f, 0x20, 0x43, 0x68, 0x65, 0x6e, 0x00, 0x00, 0x2c}
var r = bytes.NewReader(b)

type Person struct {
	Name   [10]byte
	Age    uint8
	Weight uint8
	err    error
}

func parse(r io.Reader) (*Person, error) {
	var p Person
	if err := binary.Read(r, binary.BigEndian, &p.Name); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &p.Age); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &p.Weight); err != nil {
		return nil, err
	}
	return &p, nil
}

func TestFunc(t *testing.T) {
	fmt.Println(parse(r))
}
