package ex05textbox

import "github.com/u00io/nuiforms/ui"

func Run() {
	form := ui.NewForm()
	textBox := ui.NewTextBox()
	form.Panel().AddWidget(textBox)
	form.Exec()
}
