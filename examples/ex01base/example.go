package ex01base

import "github.com/u00io/nuiforms/ui"

func NewExampleForm() *ui.Form {
	form := ui.NewForm()
	form.SetSize(400, 200)
	lbl := form.Panel().AddLabel(0, 0, "LABEL")
	lbl.SetTextAlign(ui.HAlignCenter)
	return form
}
