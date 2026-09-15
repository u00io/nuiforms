package ex04dialog

import "github.com/u00io/nuiforms/ui"

type DialogEditItem struct {
	ui.Widget

	onOK     func()
	onCancel func()
}

func NewDialogEditItem() *DialogEditItem {
	var c DialogEditItem
	c.InitWidget()

	panelEnterName := c.AddPanel(0, 0)
	panelEnterName.AddWidget(0, 0, ui.NewLabel("Enter Name:"))
	panelEnterName.AddWidget(0, 1, ui.NewTextBox())

	c.AddVSpacer(1, 0)

	panelButtons := c.AddPanel(2, 0)
	panelButtons.AddHSpacer(0, 0)
	panelButtons.AddButton(0, 1, "OK", func() { c.onOK() })
	panelButtons.AddButton(0, 2, "Cancel", func() { c.onCancel() })

	return &c
}

func (c *DialogEditItem) ExecDialog(parent *ui.Form, onOK func(), onCancel func()) {
	c.onOK = onOK
	c.onCancel = onCancel
	form := ui.NewForm()
	form.SetTitle("Edit Item Dialog")
	form.SetAllowMinimize(false)
	form.SetAllowMaximize(false)
	form.SetSize(400, 150)
	form.Panel().AddWidget(0, 0, c)
	form.ShowModal(parent)
}
