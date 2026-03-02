package ioc

import (
	"fmt"
	"testing"
)

type Widget struct {
	X, Y int
}
type Label struct {
	Widget
	Text string
}

func (label Label) Paint() {
	fmt.Printf("%p:Label.Paint(%q)\n", &label, label.Text)
}

func TestEmbed(t *testing.T) {
	label := Label{Widget{10, 10}, "State"}
	fmt.Printf("%+v\n", label)
}


