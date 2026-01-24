package ex01base

import "github.com/u00io/nuiforms/ui"

type MainForm struct {
	ui.Widget
}

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.SetLayout(`
		<label text="This is the base form example."/>
	`, &c, nil)
	return &c
}

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	form.Panel().AddWidgetOnGrid(NewMainForm(), 0, 0)
}
