package ex06font

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	label := ui.NewLabel("Click me!")
	form.Panel().AddWidgetOnGrid(label, 0, 0)
}
