package ui

func ShowAboutDialog(title string, line1 string, line2 string, line3 string, line4 string) {
	widgetToFocusAfterClose := MainForm.focusedWidget
	width := 400
	dialog := NewDialog(title, width, 200)

	panel := NewPanel()
	txtLine1 := NewLabel(line1)
	txtLine1.SetTextAlign(HAlignCenter)
	txtLine1.SetMaxWidth(width)
	panel.AddWidgetOnGrid(txtLine1, 0, 0)
	txtLine2 := NewLabel(line2)
	txtLine2.SetTextAlign(HAlignCenter)
	txtLine2.SetMaxWidth(width)
	panel.AddWidgetOnGrid(txtLine2, 1, 0)
	txtLine3 := NewLabel(line3)
	txtLine3.SetTextAlign(HAlignCenter)
	txtLine3.SetMaxWidth(width)
	panel.AddWidgetOnGrid(txtLine3, 2, 0)
	txtLine4 := NewLabel(line4)
	txtLine4.SetTextAlign(HAlignCenter)
	txtLine4.SetMaxWidth(width)
	panel.AddWidgetOnGrid(txtLine4, 3, 0)
	dialog.ContentPanel().AddWidgetOnGrid(panel, 0, 0)

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
