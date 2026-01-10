package ui

import (
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type TextBox struct {
	Widget

	cursorPosX          int
	cursorPosY          int
	selectionLeftX      int
	selectionLeftY      int
	selectionRightX     int
	selectionRightY     int
	mouseButtonPressed  bool
	cursorWidth         int
	multiline           bool
	leftAndRightPadding int

	dragingCursor bool
	// readonly      bool
	//isPassword bool

	blockUpdate bool
	emptyText   string

	padding int

	cursorVisible         bool
	skipOneCursorBlinking bool
}

type textboxModifyCommand int

const textboxModifyCommandInsertChar textboxModifyCommand = 0
const textboxModifyCommandInsertString textboxModifyCommand = 1
const textboxModifyCommandInsertReturn textboxModifyCommand = 2
const textboxModifyCommandBackspace textboxModifyCommand = 3
const textboxModifyCommandDelete textboxModifyCommand = 4
const textboxModifyCommandSetText textboxModifyCommand = 5

type TextBoxSelection struct {
	X1, Y1, X2, Y2 int
	Text           string
}

func NewTextBox() *TextBox {
	var c TextBox
	c.InitWidget()
	c.SetTypeName("TextBox")

	c.SetAutoFillBackground(true)
	c.SetElevation(-3)

	c.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		c.KeyDown(key, mods)
		return true
	})

	c.SetOnChar(func(char rune, mods nuikey.KeyModifiers) bool {
		c.KeyChar(char, mods)
		return true
	})

	c.SetOnPaint(func(cnv *Canvas) {
		c.Draw(cnv, c.innerWidth, c.innerHeight)
	})

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		c.MouseDown(button, x, y, mods)
		return true
	})

	c.SetOnMouseMove(func(x, y int, mods nuikey.KeyModifiers) bool {
		c.MouseMove(x, y, mods)
		return true
	})

	c.SetOnMouseUp(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		c.MouseUp(button, x, y, mods)
		return true
	})

	c.AddTimer(250, func() {
		c.timerCursorBlinking()
	})

	c.SetCanBeFocused(true)
	c.multiline = false
	c.SetXExpandable(true)
	c.SetYExpandable(false)
	c.SetMinSize(100, DefaultUiLineHeight)
	c.SetMaxSize(10000, DefaultUiLineHeight)

	//c.lines = make([]string, 1)
	c.cursorWidth = 1
	c.leftAndRightPadding = 0
	c.multiline = false
	c.cursorVisible = true
	c.ScrollToBegin()
	c.updateInnerSize()
	c.emptyText = "Type here..."

	c.padding = 4

	return &c
}

func (c *TextBox) Lines() []string {
	text := c.Text()
	lines := strings.Split(strings.Replace(text, "\r", "", -1), "\n")
	return lines
}

func (c *TextBox) SetReadOnly(readonly bool) {
	c.SetProp("readonly", readonly)
	UpdateMainForm()
}

func (c *TextBox) SetIsPassword(isPassword bool) {
	c.SetProp("isPassword", isPassword)
	UpdateMainForm()
}

func (c *TextBox) SetOnTextChanged(onTextChanged func()) {
	c.SetPropFunction("ontextchanged", onTextChanged)
}

type EventTextboxKeyDown struct {
	Key       nuikey.Key
	Mods      nuikey.KeyModifiers
	Processed bool
}

func (c *TextBox) SetOnTextBoxKeyDown(onKeyDown func()) {
	c.SetPropFunction("onkeydown", onKeyDown)
}

func (c *TextBox) timerCursorBlinking() {
	if MainForm.focusedWidget != nil {
		if MainForm.focusedWidget.Id() == c.id {
			if !c.skipOneCursorBlinking {
				c.cursorVisible = !c.cursorVisible
				UpdateMainForm()
			}
			c.skipOneCursorBlinking = false
		}
	}
}

func (c *TextBox) redraw() {
}

func (c *TextBox) SetProp(key string, value any) {
	c.Widget.SetProp(key, value)
	if key == "text" {
		c.setText(c.GetPropString("text", ""), false)
	}
}

func (c *TextBox) setText(text string, updateProp bool) {
	c.redraw()
	var modifiers nuikey.KeyModifiers
	c.modifyText(textboxModifyCommandSetText, modifiers, text)
	c.updateInnerSize()
	c.ScrollToBegin()
	UpdateMainForm()

	if updateProp {
		c.SetProp("text", text)
	}
}

func (c *TextBox) SetText(text string) {
	c.setText(text, true)
}

func (c *TextBox) SetEmptyText(text string) {
	c.redraw()
	c.emptyText = text
	c.updateInnerSize()
	c.ScrollToBegin()
	UpdateMainForm()
}

func (c *TextBox) Text() string {
	return c.GetPropString("text", "")
}

func (c *TextBox) SetMultiline(multiline bool, w *Widget) {
	c.multiline = multiline
	if c.multiline {
		w.allowScrollX = true
		w.allowScrollY = true
		w.SetXExpandable(true)
		w.SetYExpandable(true)
		//c.verticalScrollVisible.SetOwnValue(true)
		//c.horizontalScrollVisible.SetOwnValue(true)
	} else {
		w.SetXExpandable(true)
		w.SetYExpandable(false)
	}
	c.updateInnerSize()
	UpdateMainForm()
}

func (c *TextBox) AssemblyText(lines []string) string {
	result := ""
	for pos, line := range lines {
		result += line
		if pos < len(lines)-1 {
			result += "\n"
		}
	}
	return result
}

func (c *TextBox) updateInnerSize() {
	lines := c.Lines()

	_, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), "0")
	if err != nil {
		return
	}
	c.innerHeight = textHeight * len(lines)

	var maxTextWidth int
	for _, line := range lines {
		textWidth, _, err := MeasureText(c.FontFamily(), c.FontSize(), line)
		if err != nil {
			return
		}
		if textWidth > maxTextWidth {
			maxTextWidth = textWidth
		}
	}
	c.innerWidth = maxTextWidth + c.leftAndRightPadding*3
	if c.multiline {
		c.allowScrollY = true
	}

	if !c.multiline {
		c.innerHeight = c.Height()
	}
}

func (c *TextBox) lineToPasswordChars(line string) string {
	if c.GetPropBool("isPassword", false) {
		lenOfLine := utf8.RuneCountInString(line)
		line = ""
		for i := 0; i < lenOfLine; i++ {
			line += "*"
		}
	}
	return line
}

func (c *TextBox) Draw(ctx *Canvas, width, height int) {
	lines := c.Lines()

	oneLineHeight := c.OneLineHeight()

	var yStaticOffset int
	if c.multiline {
		yStaticOffset = 1
	} else {
		yStaticOffset = (c.Height() - oneLineHeight) / 2
		//yStaticOffset = 0
	}

	_ = yStaticOffset

	// Selection
	if len(c.selectedLines()) > 0 {
		selection := c.selectionRange()
		for selY := selection.Y1; selY <= selection.Y2; selY++ {
			lineCharPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), lines[selY])

			if err != nil {
				return
			}
			for i := 0; i < len(lineCharPos); i++ {
				lineCharPos[i] = lineCharPos[i] + c.leftAndRightPadding
			}

			selXBegin := 0
			selXWidth := lineCharPos[len(lineCharPos)-1]
			if selY == selection.Y1 {
				selXBegin = lineCharPos[selection.X1]
				selXWidth = lineCharPos[len(lineCharPos)-1] - selXBegin
			}
			if selY == selection.Y2 {
				if selection.X2 < len(lineCharPos) {
					selXWidth = lineCharPos[selection.X2] - selXBegin
				}
			}

			rectY := selY * oneLineHeight

			if !c.multiline {
				rectY = yStaticOffset
			}

			ctx.FillRect(selXBegin, rectY, selXWidth, oneLineHeight, c.BackgroundColorWithAddElevation(10))
		}
	}

	// Text
	yOffset := 0

	for _, line := range lines {
		line = c.lineToPasswordChars(line)
		//ctx.SetColor(color.RGBA{0x88, 0x88, 0x88, 0xff}) // c.foregroundColor.Color()
		ctx.SetColor(c.ForegroundColor())
		_, textHeightInLine, err := MeasureText(c.FontFamily(), c.FontSize(), line)
		ctx.SetHAlign(HAlignLeft)
		ctx.SetVAlign(VAlignCenter)
		ctx.SetFontFamily(c.FontFamily())
		ctx.SetFontSize(c.FontSize())
		ctx.DrawText(c.leftAndRightPadding, yStaticOffset+yOffset, width-c.leftAndRightPadding*2, textHeightInLine, line)

		if err != nil {
			return
		}
		yOffset += oneLineHeight
	}

	focus := MainForm.focusedWidget == c

	// Cursor
	if focus && c.cursorVisible {
		charPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), c.lineToPasswordChars(lines[c.cursorPosY]))
		for i := 0; i < len(charPos); i++ {
			charPos[i] = charPos[i] + c.leftAndRightPadding
		}
		if err != nil {
			return
		}
		cursorPosInPixels := charPos[c.cursorPosX]
		curX := cursorPosInPixels - (c.cursorWidth / 2)
		curY := yStaticOffset + c.cursorPosY*oneLineHeight
		ctx.FillRect(curX, curY, c.cursorWidth, oneLineHeight, color.RGBA{0x00, 0xFF, 0x00, 0xFF}) // c.foregroundColor.Color())
	}

	if c.Text() == "" && c.emptyText != "" && !focus {
		ctx.SetHAlign(HAlignLeft)
		ctx.SetVAlign(VAlignCenter)
		ctx.DrawText(c.leftAndRightPadding, 0, c.w, c.h, c.emptyText)
	}
}

func (c *TextBox) KeyChar(ch rune, mods nuikey.KeyModifiers) {
	if c.GetPropBool("readonly", false) {
		return
	}

	c.redraw()
	if ch < 32 {
		return
	}

	c.modifyText(textboxModifyCommandInsertChar, mods, ch)
}

func (c *TextBox) cutSelected() {
	if len(c.selectedLines()) == 0 {
		return
	}
	selectedText := c.SelectedText()
	if selectedText == "" {
		return
	}
	ClipboardSetText(selectedText)
	c.modifyText(textboxModifyCommandDelete, nuikey.KeyModifiers{}, nil)
}

func (c *TextBox) copySelected() {
	if len(c.selectedLines()) == 0 {
		return
	}

	selectedText := c.SelectedText()

	if selectedText == "" {
		return
	}

	ClipboardSetText(selectedText)
}

func (c *TextBox) paste() {
	if c.GetPropBool("readonly", false) {
		return
	}

	text, err := ClipboardGetText()
	if err != nil {
		return
	}

	if text == "" {
		return
	}

	c.modifyText(textboxModifyCommandInsertString, nuikey.KeyModifiers{}, text)
}

func (c *TextBox) KeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {

	keyDownFunc := c.GetPropFunction("onkeydown")
	if keyDownFunc != nil {
		var ev EventTextboxKeyDown
		ev.Key = key
		ev.Mods = mods
		ev.Processed = false
		PushEvent(&ev)
		keyDownFunc()
		PopEvent()
		if ev.Processed {
			return true
		}
	}

	c.redraw()

	if mods.Ctrl && key == nuikey.KeyA {
		c.SelectAllText()
		return true
	}

	if mods.Ctrl && key == nuikey.KeyX {
		c.cutSelected()
		return true
	}

	if mods.Ctrl && key == nuikey.KeyV {
		c.paste()
		return true
	}

	if mods.Ctrl && key == nuikey.KeyC {
		c.copySelected()
		return true
	}

	if key == nuikey.KeyArrowLeft {
		c.moveCursor(c.cursorPosX-1, c.cursorPosY, mods)
		return true
	}

	if key == nuikey.KeyArrowRight {
		c.moveCursor(c.cursorPosX+1, c.cursorPosY, mods)
		return true
	}

	if key == nuikey.KeyArrowUp {
		c.moveCursor(c.cursorPosX, c.cursorPosY-1, mods)
		return true
	}

	if key == nuikey.KeyArrowDown {
		c.moveCursor(c.cursorPosX, c.cursorPosY+1, mods)
		return true
	}

	if key == nuikey.KeyHome {
		c.moveCursor(0, c.cursorPosY, mods)
		return true
	}

	if key == nuikey.KeyEnter {
		if c.GetPropBool("readonly", false) {
			return false
		}
		return c.insertReturn(mods)
	}

	if key == nuikey.KeyEnd {
		lines := c.Lines()
		if c.cursorPosY >= len(lines) {
			return true
		}
		runes := []rune(lines[c.cursorPosY])
		c.moveCursor(len(runes), c.cursorPosY, mods)
		return true
	}

	if key == nuikey.KeyBackspace {
		if !c.GetPropBool("readonly", false) {
			c.modifyText(textboxModifyCommandBackspace, mods, nil)
		}
		return true
	}

	if key == nuikey.KeyDelete {
		if !c.GetPropBool("readonly", false) {
			c.modifyText(textboxModifyCommandDelete, mods, nil)
		}
		return true
	}

	return false
}

func (c *TextBox) KeyUp(key nuikey.Key, mods nuikey.KeyModifiers) {
}

func (c *TextBox) MouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	if button == nuimouse.MouseButtonLeft {
		c.redraw()
		c.mouseButtonPressed = true
		c.moveCursorNearPoint(x, y, mods)
		c.selectionLeftX = c.cursorPosX
		c.selectionLeftY = c.cursorPosY
		c.selectionRightX = c.cursorPosX
		c.selectionRightY = c.cursorPosY
		c.dragingCursor = true
		c.cursorVisible = true
		c.skipOneCursorBlinking = true
		UpdateMainForm()
	}
}

func (c *TextBox) MouseMove(x int, y int, mods nuikey.KeyModifiers) {
	c.redraw()
	if c.mouseButtonPressed {
		c.moveCursorNearPoint(x, y, mods)
	}
	UpdateMainForm()
}

func (c *TextBox) moveCursorNearPoint(x, y int, modifiers nuikey.KeyModifiers) {
	lines := c.Lines()

	_, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), "0")
	if err != nil {
		return
	}
	lineNumber := y / textHeight

	if lineNumber >= len(lines) {
		lineNumber = len(lines) - 1
	}

	if lineNumber < 0 {
		lineNumber = 0
	}

	charPos, _ := GetCharPositions(c.FontFamily(), c.FontSize(), lines[lineNumber])
	for i := 0; i < len(charPos); i++ {
		charPos[i] = charPos[i] + c.leftAndRightPadding
	}

	if len(charPos) == 1 {
		c.moveCursor(0, lineNumber, modifiers)
		return
	}

	if x < charPos[1]-(charPos[1]-charPos[0])/2 {
		c.moveCursor(0, lineNumber, modifiers)
	}

	for pos := 1; pos < len(charPos)-1; pos++ {
		left := charPos[pos] - (charPos[pos]-charPos[pos-1])/2
		right := charPos[pos] + (charPos[pos+1]-charPos[pos])/2
		if x >= left && x < right {
			c.moveCursor(pos, lineNumber, modifiers)
			break
		}
	}

	if x > charPos[len(charPos)-1] {
		c.moveCursor(len(charPos)-1, lineNumber, modifiers)
	}
}

func (c *TextBox) MouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	c.dragingCursor = false
	c.redraw()
	c.mouseButtonPressed = false
	UpdateMainForm()
}

func (c *TextBox) insertReturn(modifiers nuikey.KeyModifiers) bool {
	if !c.multiline {
		return false
	}

	c.modifyText(textboxModifyCommandInsertReturn, modifiers, nil)
	return true
}

func (c *TextBox) selectionRange() TextBoxSelection {
	var result TextBoxSelection
	//var res1X, res1Y, res2X, res2Y int
	if c.selectionLeftY > c.selectionRightY {
		result.Y1 = c.selectionRightY
		result.Y2 = c.selectionLeftY
		result.X1 = c.selectionRightX
		result.X2 = c.selectionLeftX
	}

	if c.selectionLeftY < c.selectionRightY {
		result.Y2 = c.selectionRightY
		result.Y1 = c.selectionLeftY
		result.X2 = c.selectionRightX
		result.X1 = c.selectionLeftX
	}

	if c.selectionLeftY == c.selectionRightY {
		result.Y1 = c.selectionLeftY
		result.Y2 = c.selectionRightY

		if c.selectionLeftX > c.selectionRightX {
			result.X1 = c.selectionRightX
			result.X2 = c.selectionLeftX
		} else {
			result.X2 = c.selectionRightX
			result.X1 = c.selectionLeftX
		}
	}

	return result
}

func (c *TextBox) selectedLines() []int {
	var result []int
	result = make([]int, 0)
	selection := c.selectionRange()
	if selection.Y2 != selection.Y1 {
		for i := selection.Y1; i <= selection.Y2; i++ {
			result = append(result, i)
		}
	} else {
		if selection.X1 != selection.X2 {
			result = append(result, selection.Y1)
		}
	}
	return result
}

func (c *TextBox) moveCursor(posX int, posY int, modifiers nuikey.KeyModifiers) {
	lines := c.Lines()

	if posY < 0 {
		return
	}

	if posY >= len(lines) {
		return
	}

	runes := []rune(lines[posY])

	if posX < 0 {
		return
	}

	if posX > len(runes) {
		posX = len(runes)
	}

	c.cursorPosX = posX
	c.cursorPosY = posY

	if !modifiers.Shift && !c.mouseButtonPressed {
		c.clearSelection()
	}

	if modifiers.Shift || c.dragingCursor {
		c.selectionRightX = c.cursorPosX
		c.selectionRightY = c.cursorPosY
	}

	if !c.blockUpdate {
		c.ensureVisibleCursor()
	}
	UpdateMainForm()
}

func (c *TextBox) SelectedText() string {
	lines := c.Lines()
	result := ""

	//lines := make([]string, 0)
	selection := c.selectionRange()

	if selection.Y1 == selection.Y2 {
		runes1 := []rune(lines[selection.Y1])
		result += string(runes1[selection.X1:selection.X2])
	} else {
		runes1 := []rune(lines[selection.Y1])
		result += string(runes1[selection.X1:])
		result += "\n"

		if selection.Y2-selection.Y1 > 1 {
			for row := selection.Y1 + 1; row < selection.Y2; row++ {
				result += lines[row]
				result += "\n"
			}
		}

		runes2 := []rune(lines[selection.Y2])
		result += string(runes2[0:selection.X2])
	}

	return result
}

func (c *TextBox) removeSelectedText(modifiers nuikey.KeyModifiers) (bool, []string, int, int) {
	oldLines := c.Lines()
	lines := make([]string, 0)
	modified := false
	selection := c.selectionRange()
	curPosX := c.cursorPosX
	curPosY := c.cursorPosY
	if len(c.selectedLines()) > 0 {
		lines = append(lines, oldLines[0:selection.Y1]...)
		runes1 := []rune(oldLines[selection.Y1])
		runes2 := []rune(oldLines[selection.Y2])
		lines = append(lines, string(runes1[0:selection.X1])+string(runes2[selection.X2:]))
		lines = append(lines, oldLines[selection.Y2+1:]...)
		modified = true
		curPosX = selection.X1
		curPosY = selection.Y1
	} else {
		lines = append(lines, oldLines...)
	}

	return modified, lines, curPosX, curPosY
}

func (c *TextBox) ensureVisibleCursor() {
	lines := c.Lines()
	_, oneLineHeight, _ := MeasureText(c.FontFamily(), c.FontSize(), "Q")
	charPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), lines[c.cursorPosY])
	for i := 0; i < len(charPos); i++ {
		charPos[i] = charPos[i] + c.leftAndRightPadding
	}
	if err != nil {
		return
	}
	cursorPosInPixels := charPos[c.cursorPosX]
	curX := cursorPosInPixels - (c.cursorWidth / 2)
	curY := c.cursorPosY * oneLineHeight
	// ctx.FillRect(curX, curY, c.cursorWidth, oneLineHeight)
	//c.ScrollEnsureVisible(curX, curY)
	//c.ScrollEnsureVisible(curX+c.cursorWidth, curY+oneLineHeight)
	c.ScrollEnsureVisible(curX, curY)
	c.ScrollEnsureVisible(curX+c.cursorWidth, curY+oneLineHeight)
}

func (c *TextBox) clearSelection() {
	c.selectionLeftX = c.cursorPosX
	c.selectionLeftY = c.cursorPosY
	c.selectionRightX = c.cursorPosX
	c.selectionRightY = c.cursorPosY
}

func (c *TextBox) modifyText(cmd textboxModifyCommand, modifiers nuikey.KeyModifiers, data interface{}) {
	c.redraw()

	valid := true
	selectedTextRemoved, lines, curPosX, curPosY := c.removeSelectedText(modifiers)

	switch cmd {
	case textboxModifyCommandInsertChar:
		{
			out := []rune(lines[curPosY])
			left := string(out[0:curPosX])
			right := string(out[curPosX:])
			lines[curPosY] = left + string(data.(rune)) + right
			curPosX += 1
		}
	case textboxModifyCommandInsertReturn:
		{
			runes := []rune(lines[curPosY])
			left := string(runes[0:curPosX])
			right := string(runes[curPosX:])
			linesBefore := lines[0:curPosY]
			linesAfter := lines[curPosY:]
			lines = append(linesBefore, right)
			lines = append(lines, linesAfter...)
			lines[curPosY] = left
			curPosX = 0
			curPosY++
		}
	case textboxModifyCommandBackspace:
		{
			runes := []rune(lines[curPosY])
			if !selectedTextRemoved {
				if curPosX > 0 {
					left := string(runes[0 : curPosX-1])
					right := string(runes[curPosX:])
					lines[curPosY] = left + right
					curPosX = curPosX - 1
				} else {
					if curPosY > 0 {
						runes := []rune(lines[curPosY-1])
						newCursorPosX := len(runes)
						linesTemp := make([]string, 0)
						linesTemp = append(linesTemp, lines[0:curPosY]...)
						linesTemp[curPosY-1] += lines[curPosY]
						linesTemp = append(linesTemp, lines[curPosY+1:]...)
						lines = linesTemp
						curPosX = newCursorPosX
						curPosY = curPosY - 1
					}
				}
			}
		}
	case textboxModifyCommandDelete:
		{
			runes := []rune(lines[curPosY])
			if !selectedTextRemoved {
				if curPosX < len(runes) {
					left := string(runes[0:curPosX])
					right := string(runes[curPosX+1:])
					lines[curPosY] = left + right
				} else {
					if curPosY < len(lines)-1 {
						linesTemp := make([]string, 0)
						linesTemp = append(linesTemp, lines[0:curPosY+1]...)
						linesTemp[curPosY] += lines[curPosY+1]
						linesTemp = append(linesTemp, lines[curPosY+2:]...)
						lines = linesTemp
					}
				}
			}
		}
	case textboxModifyCommandSetText:
		{
			lines = strings.Split(strings.Replace(data.(string), "\r", "", -1), "\n")
			//curPosX = 0
			//curPosY = 0
		}
	case textboxModifyCommandInsertString:
		{
			c.blockUpdate = true
			runes := string(data.(string))
			for _, ch := range runes {
				if ch < 32 {
					if ch == 10 {
						c.insertReturn(modifiers)
					}
				}

				c.KeyChar(ch, modifiers)
			}
			lines = c.Lines()
			curPosX = c.cursorPosX
			curPosY = c.cursorPosY
			c.blockUpdate = false
		}
	}

	if valid {
		newText := c.AssemblyText(lines)
		c.Widget.SetProp("text", newText)
		c.updateInnerSize()
		c.moveCursor(curPosX, curPosY, modifiers)

		if !c.blockUpdate {
			c.clearSelection()
			c.updateInnerSize()

			f := c.GetPropFunction("ontextchanged")
			if f != nil {
				PushEvent(nil)
				f()
				PopEvent()
			}
		}

	}

	UpdateMainForm()
}

func (c *TextBox) SelectAllText() {
	lines := c.Lines()
	runesLast := []rune(lines[len(lines)-1])
	c.selectionLeftX = 0
	c.selectionLeftY = 0
	c.selectionRightX = len(runesLast)
	c.selectionRightY = len(lines) - 1
}

func (c *TextBox) MoveCursorToEnd() {
	lines := c.Lines()
	runes := []rune(lines[c.cursorPosY])
	c.moveCursor(len(runes), c.cursorPosY, nuikey.KeyModifiers{})
}

func (c *TextBox) ScrollToBegin() {
	c.ScrollEnsureVisible(0, 0)
	c.ScrollEnsureVisible(0, 1)
	//c.ScrollEnsureVisible(0, 0)
	//c.ScrollEnsureVisible(0, 1)
}

func (c *TextBox) OneLineHeight() int {
	_, fontHeight, _ := MeasureText(c.FontFamily(), c.FontSize(), "1Qg")
	return fontHeight
}

/*
func (c *TextBox) MinHeight() int {
	return c.OneLineHeight() + 4 + c.TopBorderWidth() + c.BottomBorderWidth()
}

func (c *TextBox) AcceptsReturn() bool {
	return c.multiline
}
*/

/*func (c *TextBox) FontFamily() string {
	return "robotomono"
}

func (c *TextBox) FontSize() float64 {
	return 16
}*/
