package err

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

type Reader struct {
	r   io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func parseStruct(input io.Reader) (*Person, error) {
	var p Person
	r := Reader{r: input}

	r.read(&p.Name)
	r.read(&p.Age)
	r.read(&p.Weight)

	if r.err != nil {
		return nil, r.err
	}

	return &p, nil
}

func TestStruct(t *testing.T) {
	fmt.Println(parseStruct(r))
}
