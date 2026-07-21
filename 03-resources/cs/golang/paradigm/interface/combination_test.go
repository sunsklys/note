package _interface

import (
	"fmt"
	"testing"
)

type Printable interface {
	PrintStr()
}

type WithName struct {
	Name string
}

func (w WithName) PrintStr() {
	fmt.Println(w.Name)
}

type CountryC struct {
	WithName
}

type CityC struct {
	WithName
}

func TestCombination(t *testing.T) {
	c1 := CountryC{WithName{"China"}}
	c2 := CityC{WithName{"Beijing"}}
	c1.PrintStr()
	c2.PrintStr()
}
