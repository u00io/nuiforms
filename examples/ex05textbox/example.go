package ex05textbox

import "github.com/u00io/nuiforms/ui"

func Run() {
	form := ui.NewForm()
	textBox := ui.NewTextBox()
	textBox.SetPosition(50, 50)
	textBox.SetSize(200, 30)
	form.Panel().AddWidget(textBox)
	form.Exec()
}
