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
<frame padding="50">
	<frame padding="50">
		<button />
	</frame>
</frame>
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
