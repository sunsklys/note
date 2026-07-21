package ioc

import (
	"fmt"
	"testing"
)

type Button struct {
	Label
}

func NewButton(x, y int, text string) Button {
	return Button{
		Label: Label{
			Widget: Widget{
				X: x,
				Y: y,
			},
			Text: text,
		},
	}
}

func (button Button) Paint() {
	fmt.Printf("Button.Paint(%s)\n", button.Text)
}
func (button Button) Click() {
	fmt.Printf("Button.Click(%s)\n", button.Text)
}

type ListBox struct {
	Widget
	Texts []string
	Index int
}

func (listBox ListBox) Paint() {
	fmt.Printf("ListBox.Paint(%q)\n", listBox.Texts)
}
func (listBox ListBox) Click() {
	fmt.Printf("ListBox.Click(%q)\n", listBox.Texts)
}

type Painter interface {
	Paint()
}

type Clicker interface {
	Click()
}

func TestReWrite(t *testing.T) {
	label := Label{
		Widget: Widget{
			X: 1,
			Y: 2,
		},
		Text: "abc",
	}
	button1 := Button{Label{Widget{10, 70}, "OK"}}
	button2 := NewButton(50, 70, "Cancel")
	listBox := ListBox{Widget{10, 40},
		[]string{"AL", "AK", "AZ", "AR"}, 0}

	for _, painter := range []Painter{label, listBox, button1, button2} {
		painter.Paint()
	}

	fmt.Println()
	for _, widget := range []any{label, listBox, button1, button2} {
		widget.(Painter).Paint()
		if clicker, ok := widget.(Clicker); ok {
			clicker.Click()
		}
	}
}
