package ex01base

import "github.com/u00io/nuiforms/ui"

func NewExampleForm() *ui.Form {
	form := ui.NewForm()
	form.Panel().AddWidget(ui.NewLabel("123"), 0, 0)
	return form
}

func Run(form *ui.Form) {
	f := NewExampleForm()
	f.SetSize(400, 200)
	f.ShowModal(form)
}
