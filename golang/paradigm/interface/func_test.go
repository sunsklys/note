package _interface

import (
	"fmt"
	"testing"
)

type Person struct {
	Name   string
	Age    int
	Sexual string
}

func (p *Person) Print() {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
		p.Name, p.Sexual, p.Age)
}

func PrintPerson(p *Person) {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
		p.Name, p.Sexual, p.Age)
}

func TestPrintFunc(t *testing.T) {
	var p = Person{
		Name:   "Hao Chen",
		Sexual: "Male",
		Age:    44,
	}

	PrintPerson(&p)
	p.Print()
}
