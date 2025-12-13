package ui

func ShowEnterStringDialog(title string, messageText string, initialValue string, onSubmit func(value string)) {
	widgetToFocusAfterClose := MainForm.focusedWidget
	width := 400
	dialog := NewDialog(title, width, 200)
	txtMessage := NewLabel(messageText)
	txtMessage.SetTextAlign(HAlignCenter)
	txtMessage.SetMaxWidth(width)
	dialog.ContentPanel().AddWidgetOnGrid(txtMessage, 0, 0)

	txtValue := NewTextBox()
	txtValue.SetText(initialValue)
	dialog.ContentPanel().AddWidgetOnGrid(txtValue, 1, 0)

	dialog.ContentPanel().AddWidgetOnGrid(NewVSpacer(), 2, 0)

	panelButtons := NewPanel()
	panelButtons.AddWidgetOnGrid(NewHSpacer(), 0, 0)
	btnOK := NewButton("OK")
	btnOK.SetOnButtonClick(func() {
		if onSubmit != nil {
			onSubmit(txtValue.Text())
		}
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidgetOnGrid(btnOK, 0, 1)

	btnCancel := NewButton("Cancel")
	btnCancel.SetOnButtonClick(func() {
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidgetOnGrid(btnCancel, 0, 2)

	dialog.ContentPanel().AddWidgetOnGrid(panelButtons, 3, 0)
	dialog.ShowDialog()
	btnOK.Focus()
}
