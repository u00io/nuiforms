package ui

import "strings"

func ShowMessageBox(title string, messageText string) {
	width := 500
	charsPerLine := 20

	smartWordWrap := func(text string, maxChars int) []string {
		words := strings.Fields(text)
		lines := make([]string, 0)
		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word)+1 > maxChars {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = ""
			}
			currentLine += word + " "
		}
		if len(currentLine) > 0 {
			lines = append(lines, strings.TrimSpace(currentLine))
		}
		return lines
	}

	// wrap the text to fit within the message box width
	lines1 := smartWordWrap(messageText, charsPerLine)
	lines := make([]string, 0)
	/*if len(lines1) >= 5 {
		lines = lines1[:5]
		lines = append(lines, "...")
	}*/
	_ = lines1

	lines = append(lines, "11111")
	lines = append(lines, "22222")
	lines = append(lines, "33333")
	lines = append(lines, "44444")
	lines = append(lines, "55555")

	dialogHeight := 140 + len(lines)*(30+10)

	widgetToFocusAfterClose := MainForm.focusedWidget
	dialog := NewDialog(title, width, dialogHeight)
	dialog.ContentPanel().SetLayout(`
		<column>
			<column id="colLines"/>
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
