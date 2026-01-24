package ex15dialog

import (
	"github.com/u00io/nuiforms/ui"
)

type DialogEnterName struct {
	dialog       *ui.Dialog
	txtName      *ui.TextBox
	panelContent *ui.Panel
	panelButtons *ui.Panel
	btnOk        *ui.Button
	btnCancel    *ui.Button
}

func NewDialogEnterName() *DialogEnterName {
	var c DialogEnterName
	c.dialog = ui.NewDialog("Enter Name", 300, 200)
	c.panelContent = ui.NewPanel()
	c.txtName = ui.NewTextBox()
	c.panelContent.AddWidgetOnGrid(ui.NewLabel("Name:"), 0, 0)
	c.panelContent.AddWidgetOnGrid(c.txtName, 0, 1)
	c.panelContent.AddWidgetOnGrid(ui.NewVSpacer(), 1, 0)
	c.dialog.ContentPanel().AddWidgetOnGrid(c.panelContent, 0, 0)
	c.dialog.SetCloseByClickOutside(false)

	c.btnOk = ui.NewButton("OK")
	c.btnCancel = ui.NewButton("Cancel")
	c.panelButtons = ui.NewPanel()
	c.panelButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 0)
	c.panelButtons.AddWidgetOnGrid(c.btnOk, 0, 1)
	c.panelButtons.AddWidgetOnGrid(c.btnCancel, 0, 2)
	//c.panelButtons.SetBackgroundColor(color.RGBA{R: 200, G: 20, B: 200, A: 255})
	c.dialog.ContentPanel().AddWidgetOnGrid(c.panelButtons, 1, 0)

	c.dialog.SetAcceptButton(c.btnOk)
	c.dialog.SetRejectButton(c.btnCancel)
	return &c
}

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	btn := ui.NewButton("Click me")
	btn.SetOnClick(func() {
		dialog := NewDialogEnterName()
		dialog.dialog.ShowDialog()
		dialog.dialog.OnAccept = func() {
			name := dialog.txtName.Text()
			form.SetTitle("Hello " + name)
		}
		dialog.dialog.OnReject = func() {
			form.SetTitle("Dialog was cancelled")
		}
	})
	form.Panel().AddWidgetOnGrid(btn, 0, 0)
	form.Panel().AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)
	form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), 10, 0)
}
