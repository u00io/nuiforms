package ui

func ShowMessageBox(title string, messageText string) {
	width := 500

	fontFamily := ThemeFontFamily()
	fontSize := ThemeFontSize()

	charsPerLine := 20

	oneSymbolWidth, _, err := MeasureText(fontFamily, fontSize, "W")
	if err == nil {
		charsPerLine = (width - 20) / (oneSymbolWidth)
	}

	// wrap the text to fit within the message box width
	lines := MakeLinesFromStringWithWordWrapping(messageText, charsPerLine)
	if len(lines) >= 5 {
		lines = lines[:5]
		lines = append(lines, "...")
	}

	lineHeight := DefaultUiLineHeight + 2

	dialogHeight := 100 + len(lines)*lineHeight

	widgetToFocusAfterClose := MainForm.focusedWidget
	dialog := NewDialog(title, width, dialogHeight)
	dialog.ContentPanel().SetLayout(`
		<column>
			<column id="colLines" spacing="2"/>
			<vspacer />
			<row>
				<hspacer />
				<button id="btnOK" text="OK" onclick="OnOKClick"/>
				<hspacer />
			</row>
		</column>
	`, map[string]func(){
		"OnOKClick": func() {
			dialog.Close()
			widgetToFocusAfterClose.Focus()
		},
	}, nil)

	colLines := dialog.FindWidgetByName("colLines").(*Panel)
	for i, line := range lines {
		lblLine := NewLabel(line)
		lblLine.SetTextAlign(HAlignLeft)
		lblLine.SetXExpandable(true)
		colLines.AddWidgetOnGrid(lblLine, i, 0)
	}
	UpdateMainFormLayout()

	dialog.ShowDialog()
	btnOK := dialog.FindWidgetByName("btnOK").(*Button)
	btnOK.Focus()
}
