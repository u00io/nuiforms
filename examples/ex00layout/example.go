package ex00layout

import (
	"github.com/u00io/nuiforms/ui"
)

type InnerWidget struct {
	ui.Widget
}

func NewInnerWidget() *InnerWidget {
	var c InnerWidget
	c.InitWidget()
	c.SetLayout(`
<column>
	<label text="This is an inner widget." />
	<vspacer />
	<label text="It is placed inside the main layout." />
	<button text="Inner Button" onclick="ClickInnerButton" />
</column>
	`, &c, nil)
	return &c
}

func (c *InnerWidget) ClickInnerButton() {
	println("Inner button clicked!")
}

type LayoutExample struct {
	ui.Widget
}

func NewLayoutExample() *LayoutExample {
	var c LayoutExample
	c.InitWidget()
	c.SetLayout(`
<column>
    <label text="Do you want to destroy the system?" />
	<vspacer />
	<widget id="InnerWidget" />
	<row>
		<hspacer />
		<button text="Yes" onclick="ClickYes" />
		<button text="No" />
    </row>
</column>
	`, &c, map[string]ui.Widgeter{
		"InnerWidget": NewInnerWidget(),
	})
	return &c
}

func (c *LayoutExample) ClickYes() {
	println("System destroyed!")
}

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Layout")
	form.Panel().AddWidgetOnGrid(NewLayoutExample(), 0, 0)
}
