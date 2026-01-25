package ex01base

import "github.com/u00io/nuiforms/ui"

type MainWidget struct {
	ui.Widget
}

func NewMainWidget() *MainWidget {
	var c MainWidget
	c.InitWidget()
	c.SetLayout(`
		<label text="This is the base form example."/>
	`, &c, nil)
	return &c
}

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	form.SetMainWidget(NewMainWidget())
}
