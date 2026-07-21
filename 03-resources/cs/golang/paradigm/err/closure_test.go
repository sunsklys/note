package err

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

func parseClosure(r io.Reader) (*Person, error) {
	var p Person
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, data)
	}

	read(&p.Name)
	read(&p.Age)
	read(&p.Weight)

	if err != nil {
		return &p, err
	}

	return &p, nil
}

func TestFuncClosure(t *testing.T) {
	fmt.Println(parseClosure(r))
}
