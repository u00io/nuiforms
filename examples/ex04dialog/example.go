package ex04dialog

import "github.com/u00io/nuiforms/ui"

func NewExampleForm() *ui.Form {
	form := ui.NewForm()
	form.Panel().AddButton(0, 0, "Dialog New Item", func() {
		dialog := NewDialogNewItem()
		form.Panel().ShowDialog("New Item Dialog", 400, 150, dialog)
	})
	form.Panel().AddButton(0, 1, "Dialog Edit Item", func() {
		dialog := NewDialogEditItem()
		dialog.ExecDialog(form, func() {}, func() {})

	})
	form.Panel().AddButton(0, 2, "Dialog Settings", func() {
		dialog := NewDialogSettings()
		dialog.ExecDialog(form, func() {}, func() {})

	})
	form.Panel().AddHSpacer(0, 10)
	form.Panel().AddVSpacer(10, 0)
	return form
}
