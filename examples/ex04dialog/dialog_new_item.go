package ex04dialog

import "github.com/u00io/nuiforms/ui"

type DialogNewItem struct {
	ui.Widget

	txtBox *ui.TextBox

	onOK     func()
	onCancel func()
}

func NewDialogNewItem() *DialogNewItem {
	var c DialogNewItem
	c.InitWidget()

	panelEnterName := c.AddPanel(0, 0)
	panelEnterName.AddWidget(0, 0, ui.NewLabel("Enter Name:"))
	c.txtBox = ui.NewTextBox()
	panelEnterName.AddWidget(0, 1, c.txtBox)

	c.AddVSpacer(1, 0)

	panelButtons := c.AddPanel(2, 0)
	panelButtons.AddHSpacer(0, 0)
	panelButtons.AddButton(0, 1, "OK", func() { c.onOK() })
	panelButtons.AddButton(0, 2, "Cancel", func() { c.onCancel() })

	c.txtBox.Focus()

	return &c
}
