package ex00gallery

import "github.com/u00io/nuiforms/ui"

type ExamplePageTextBox struct {
	ui.Widget
}

func NewExamplePageTextBox() *ExamplePageTextBox {
	var c ExamplePageTextBox
	c.InitWidget()
	return &c
}
