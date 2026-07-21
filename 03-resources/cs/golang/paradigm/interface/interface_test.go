package _interface

import (
	"fmt"
	"testing"
)

type Stringable interface {
	ToString() string
}
type Country struct {
	Name string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}

type City struct {
	Name string
}

func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

func TestInterface(t *testing.T) {
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}
	PrintStr(d1)
	PrintStr(d2)

}
