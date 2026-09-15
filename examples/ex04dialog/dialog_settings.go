package ex04dialog

import "github.com/u00io/nuiforms/ui"

type DialogSettings struct {
	ui.Widget
}

func NewDialogSettings() *DialogSettings {
	var c DialogSettings
	c.InitWidget()
	return &c
}

func (c *DialogSettings) ExecDialog(parent *ui.Form, onOK func(), onCancel func()) {
	form := ui.NewForm()
	form.Panel().AddWidget(0, 0, c)
	form.ShowModal(parent)
}
