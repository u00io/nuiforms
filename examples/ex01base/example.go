package ex01base

import "github.com/u00io/nuiforms/ui"

func NewExampleForm() *ui.Form {
	form := ui.NewForm()
	form.SetSize(400, 200)
	form.Panel().AddWidget(ui.NewLabel("123"), 0, 0)
	return form
}
