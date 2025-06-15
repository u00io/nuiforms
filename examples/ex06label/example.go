package ex06label

import "github.com/u00io/nuiforms/ui"

func Run() {
	form := ui.NewForm()
	form.Panel().SetAbsolutePositioning(true)
	label := ui.NewLabel("Click me!")
	label.SetPosition(50, 50)
	label.SetSize(200, 30)
	form.Panel().AddWidget(label)
	form.Exec()
}
