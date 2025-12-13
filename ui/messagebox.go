package ui

func ShowMessageBox(title string, messageText string) {
	widgetToFocusAfterClose := MainForm.focusedWidget
	width := 400
	dialog := NewDialog(title, width, 200)
	txtMessage := NewLabel(messageText)
	txtMessage.SetTextAlign(HAlignCenter)
	txtMessage.SetMaxWidth(width)
	dialog.ContentPanel().AddWidgetOnGrid(txtMessage, 0, 0)

	dialog.ContentPanel().AddWidgetOnGrid(NewVSpacer(), 1, 0)

	panelButtons := NewPanel()

	btnOK := NewButton("OK")
	btnOK.SetOnButtonClick(func() {
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidgetOnGrid(NewHSpacer(), 0, 0)
	panelButtons.AddWidgetOnGrid(btnOK, 0, 1)
	panelButtons.AddWidgetOnGrid(NewHSpacer(), 0, 2)

	dialog.ContentPanel().AddWidgetOnGrid(panelButtons, 2, 0)
	dialog.ShowDialog()
	btnOK.Focus()
}
