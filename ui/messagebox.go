package ui

import "strings"

const (
	// messageBoxMaxLines is the maximum number of text lines shown in a message box.
	messageBoxMaxLines = 20

	messageBoxMinWidth = 250
	messageBoxMaxWidth = 600

	// messageBoxTextHPadding is the horizontal space reserved around the text
	// (panel/cell paddings, window borders, etc.) when computing the window width.
	messageBoxTextHPadding = 40

	// messageBoxChromeHeight is the vertical space taken by everything except
	// the text itself (spacer, button row, paddings).
	messageBoxChromeHeight = 120
)

func ShowMessageBox(parentWidget Widgeter, title string, text string) {
	c := NewForm()
	c.SetTitle(title)
	c.SetAllowMinimize(false)
	c.SetAllowMaximize(false)

	fontFamily := ThemeFontFamily()
	fontSize := DefaultFontSize

	maxTextWidth := messageBoxMaxWidth - messageBoxTextHPadding
	lines := wrapTextLines(text, fontFamily, fontSize, maxTextWidth)
	if len(lines) > messageBoxMaxLines {
		lines = lines[:messageBoxMaxLines]
		lines[len(lines)-1] = lines[len(lines)-1] + "..."
	}

	maxLineWidth := 0
	for _, line := range lines {
		lineWidth, _, err := MeasureText(fontFamily, fontSize, line)
		if err == nil && lineWidth > maxLineWidth {
			maxLineWidth = lineWidth
		}
	}

	_, lineHeight, err := MeasureText(fontFamily, fontSize, "Йg")
	if err != nil || lineHeight <= 0 {
		lineHeight = DefaultUiLineHeight
	}
	textHeight := lineHeight * len(lines)

	width := maxLineWidth + messageBoxTextHPadding
	if width < messageBoxMinWidth {
		width = messageBoxMinWidth
	}
	if width > messageBoxMaxWidth {
		width = messageBoxMaxWidth
	}

	height := textHeight + messageBoxChromeHeight

	c.SetSize(width, height)

	lbl := c.Panel().AddLabel(0, 0, strings.Join(lines, "\r\n"))
	lbl.SetTextAlign(HAlignCenter)
	lbl.SetMinHeight(textHeight)
	c.Panel().AddVSpacer(1, 0)
	panel := c.Panel().AddPanel(2, 0)
	panel.AddHSpacer(0, 0)
	panel.AddButton(0, 1, "OK", func() { c.Close() })
	panel.AddHSpacer(0, 2)
	if parentWidget != nil {
		c.ShowModal(parentWidget.Form())
	} else {
		c.Show()
	}
}

// wrapTextLines splits text into lines that fit within maxWidth pixels when
// rendered with the given font. Existing line breaks ("\n" or "\r\n") in the
// input are preserved as paragraph breaks.
func wrapTextLines(text string, fontFamily string, fontSize float64, maxWidth int) []string {
	var lines []string

	text = strings.ReplaceAll(text, "\r\n", "\n")
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := ""
		for _, word := range words {
			candidate := word
			if currentLine != "" {
				candidate = currentLine + " " + word
			}

			candidateWidth, _, err := MeasureText(fontFamily, fontSize, candidate)
			if err == nil && candidateWidth <= maxWidth {
				currentLine = candidate
				continue
			}

			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}

			wordWidth, _, err := MeasureText(fontFamily, fontSize, word)
			if err == nil && wordWidth > maxWidth {
				// The word alone is wider than maxWidth - break it by characters.
				currentLine = breakLongWord(word, fontFamily, fontSize, maxWidth, &lines)
			} else {
				currentLine = word
			}
		}

		if currentLine != "" {
			lines = append(lines, currentLine)
		}
	}

	return lines
}

// breakLongWord splits a single word that is wider than maxWidth into chunks,
// appending all but the last chunk to lines and returning the remaining tail.
func breakLongWord(word string, fontFamily string, fontSize float64, maxWidth int, lines *[]string) string {
	chunk := ""
	for _, r := range word {
		candidate := chunk + string(r)
		candidateWidth, _, err := MeasureText(fontFamily, fontSize, candidate)
		if err == nil && candidateWidth > maxWidth && chunk != "" {
			*lines = append(*lines, chunk)
			chunk = string(r)
			continue
		}
		chunk = candidate
	}
	return chunk
}
