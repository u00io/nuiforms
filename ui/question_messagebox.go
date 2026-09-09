package ui

func ShowQuestionMessageBox(title string, messageText string, onOk func(), onCancel func()) {
	widgetToFocusAfterClose := MainForm.focusedWidget
	width := 400
	dialog := NewDialog(title, width, 200)
	txtMessage := NewLabel(messageText)
	txtMessage.SetTextAlign(HAlignCenter)
	txtMessage.SetMaxWidth(width)
	dialog.ContentPanel().AddWidget(txtMessage, 0, 0)

	dialog.ContentPanel().AddWidget(NewVSpacer(), 1, 0)

	panelButtons := NewPanel()

	btnOK := NewButton("OK")
	btnOK.SetOnClick(func() {
		if onOk != nil {
			onOk()
		}
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidget(NewHSpacer(), 0, 0)
	panelButtons.AddWidget(btnOK, 0, 1)

	btnCancel := NewButton("Cancel")
	btnCancel.SetOnClick(func() {
		if onCancel != nil {
			onCancel()
		}
		dialog.Close()
		widgetToFocusAfterClose.Focus()
	})
	panelButtons.AddWidget(btnCancel, 0, 2)

	panelButtons.AddWidget(NewHSpacer(), 0, 3)

	dialog.ContentPanel().AddWidget(panelButtons, 2, 0)
	dialog.ShowDialog()
	btnOK.Focus()
}
