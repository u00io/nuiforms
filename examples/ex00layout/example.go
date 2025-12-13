package ex00layout

import (
	"github.com/u00io/nuiforms/ui"
)

type LayoutExample struct {
	ui.Widget
}

func NewLayoutExample() *LayoutExample {
	var c LayoutExample
	c.InitWidget()

	c.SetLayout(`
<column>
    <textbox text="This is an example of a layout with an inner widget." />
</column>
	`, &c, nil)

	return &c
}

func (c *LayoutExample) ClickYes() {
	println("System destroyed!")
}

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Layout")
	form.Panel().AddWidgetOnGrid(NewLayoutExample(), 0, 0)
}
