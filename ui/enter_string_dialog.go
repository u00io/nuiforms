package ui

func ShowEnterStringDialog(title string, messageText string, initialValue string, onSubmit func(value string)) {
	widgetToFocusAfterClose := MainForm.focusedWidget
	width := 400
	dialog := NewDialog(title, width, 200)
	txtMessage := NewLabel(messageText)
	txtMessage.SetTextAlign(HAlignCenter)
	txtMessage.SetMaxWidth(width)
	dialog.ContentPanel().AddWidget(txtMessage, 0, 0)

	txtValue := NewTextBox()
	txtValue.SetText(initialValue)
	dialog.ContentPanel().AddWidget(txtValue, 1, 0)

	dialog.ContentPanel().AddWidget(NewVSpacer(), 2, 0)

	panelButtons := NewPanel()
	panelButtons.AddWidget(NewHSpacer(), 0, 0)
	btnOK := NewButton("OK")
	btnOK.SetOnClick(func() {
		if onSubmit != nil {
			onSubmit(txtValue.Text())
		}
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidget(btnOK, 0, 1)

	btnCancel := NewButton("Cancel")
	btnCancel.SetOnClick(func() {
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidget(btnCancel, 0, 2)

	dialog.ContentPanel().AddWidget(panelButtons, 3, 0)
	dialog.ShowDialog()
	btnOK.Focus()
}
