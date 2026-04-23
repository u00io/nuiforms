package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

// NumBox is a single-line numeric input for float64 with min/max clamp,
// fixed decimals formatting and spin buttons (up/down).
type NumBox struct {
	Widget

	// editing state
	text      string
	cursorPos int // rune index in text

	mouseButtonPressed bool
	draggingSelection  bool
	selectionAnchorPos int // rune index
	selectionEndPos    int // rune index
	propIsProcessing   bool

	// value state
	value       float64
	hasValue    bool
	decimals    int
	min         float64
	max         float64
	step        float64
	editingText bool // whether text may differ from formatted value

	// visual
	padding int
}

/*
Properties:
- value: float64 - current value
- decimals: int - number of digits after decimal point
- min: float64 - min allowed value
- max: float64 - max allowed value
- step: float64 - step used for buttons (if 0, computed from decimals)
*/

func NewNumBox() *NumBox {
	var c NumBox
	c.InitWidget()
	c.SetTypeName("NumBox")
	c.SetAutoFillBackground(true)
	c.SetElevation(-3)
	c.SetCanBeFocused(true)
	c.SetXExpandable(true)
	c.SetYExpandable(false)
	c.SetMinSize(120, DefaultUiLineHeight)
	c.SetMouseCursor(nuimouse.MouseCursorIBeam)

	c.padding = 4
	c.min = math.Inf(-1)
	c.max = math.Inf(1)
	c.decimals = 2
	c.step = 0

	c.SetOnPaint(func(cnv *Canvas) { c.draw(cnv) })

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		if button != nuimouse.MouseButtonLeft {
			return false
		}
		c.Focus()
		c.onMouseDownLeft(x, y, mods)
		return true
	})

	c.SetOnMouseMove(func(x, y int, mods nuikey.KeyModifiers) bool {
		c.onMouseMove(x, y, mods)
		return true
	})

	c.SetOnMouseLeave(func() {
		// Restore default cursor when leaving widget.
		c.SetMouseCursor(nuimouse.MouseCursorIBeam)
	})

	c.SetOnMouseUp(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.onMouseUpLeft(x, y, mods)
			return true
		}
		return false
	})

	c.SetOnMouseWheel(func(deltaX, deltaY int) bool {
		_ = deltaX
		c.onMouseWheel(deltaY)
		return true
	})

	c.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		return c.onKeyDown(key, mods)
	})

	c.SetOnChar(func(ch rune, mods nuikey.KeyModifiers) bool {
		return c.onChar(ch, mods)
	})

	c.SetOnFocusLost(func() {
		// Commit when focus is lost.
		c.commitText(true)
	})

	c.SetValue(0)
	return &c
}

func (c *NumBox) SetOnValueChanged(f func()) {
	c.SetPropFunction("onvaluechanged", f)
}

// SetOnChanged is an alias for SetOnValueChanged, kept for consistency with other controls.
func (c *NumBox) SetOnChanged(f func()) {
	c.SetPropFunction("onchanged", f)
}

type EventNumBoxValueChanged struct {
	NumBox *NumBox
	Value  float64
}

func (c *NumBox) fireValueChanged() {
	var ev EventNumBoxValueChanged
	ev.NumBox = c
	ev.Value = c.value

	// Support both names for the callback.
	if f := c.GetPropFunction("onvaluechanged"); f != nil {
		PushEvent(&ev)
		f()
		PopEvent()
	}
	if f := c.GetPropFunction("onchanged"); f != nil {
		PushEvent(&ev)
		f()
		PopEvent()
	}
}

func (c *NumBox) SetDecimals(decimals int) {
	if decimals < 0 {
		decimals = 0
	}
	if decimals > 16 {
		decimals = 16
	}
	c.decimals = decimals
	c.SetProp("decimals", decimals)
	c.refreshTextFromValue()
}

func (c *NumBox) Decimals() int {
	return c.decimals
}

func (c *NumBox) SetMin(min float64) {
	c.min = min
	c.SetProp("min", min)
	c.SetValue(c.value)
}

func (c *NumBox) Min() float64 {
	return c.min
}

func (c *NumBox) SetMax(max float64) {
	c.max = max
	c.SetProp("max", max)
	c.SetValue(c.value)
}

func (c *NumBox) Max() float64 {
	return c.max
}

func (c *NumBox) SetStep(step float64) {
	if step < 0 {
		step = -step
	}
	c.step = step
	c.SetProp("step", step)
}

func (c *NumBox) Step() float64 {
	if c.step > 0 {
		return c.step
	}
	// Default: 1 at decimals=0, otherwise 10^-decimals
	if c.decimals <= 0 {
		return 1
	}
	return math.Pow10(-c.decimals)
}

func (c *NumBox) clamp(v float64) float64 {
	if v < c.min {
		return c.min
	}
	if v > c.max {
		return c.max
	}
	return v
}

func (c *NumBox) SetValue(v float64) {
	if math.IsNaN(v) {
		return
	}
	v = c.clamp(v)

	changed := !c.hasValue || v != c.value
	c.value = v
	c.hasValue = true
	// Prevent recursion: SetProp triggers ProcessPropChange("value").
	c.propIsProcessing = true
	c.SetProp("value", v)
	c.propIsProcessing = false

	c.refreshTextFromValue()

	if changed {
		c.fireValueChanged()
	}
	UpdateMainForm()
}

func (c *NumBox) Value() float64 {
	return c.value
}

func (c *NumBox) GetValue() float64 {
	return c.Value()
}

func (c *NumBox) formatValue(v float64) string {
	// Avoid "-0.00" etc.
	if math.Abs(v) < 0.5*math.Pow10(-c.decimals) {
		v = 0
	}
	return fmt.Sprintf("%.*f", c.decimals, v)
}

func (c *NumBox) refreshTextFromValue() {
	c.text = c.formatValue(c.value)
	c.cursorPos = len([]rune(c.text))
	c.selectionAnchorPos = c.cursorPos
	c.selectionEndPos = c.cursorPos
	c.draggingSelection = false
	c.mouseButtonPressed = false
	c.editingText = false
}

func (c *NumBox) Text() string {
	return c.text
}

func (c *NumBox) SetText(text string) {
	c.text = text
	c.cursorPos = len([]rune(c.text))
	c.selectionAnchorPos = c.cursorPos
	c.selectionEndPos = c.cursorPos
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) hasSelection() bool {
	return c.selectionAnchorPos != c.selectionEndPos
}

func (c *NumBox) selectionRange() (from, to int) {
	from = c.selectionAnchorPos
	to = c.selectionEndPos
	if from > to {
		from, to = to, from
	}
	if from < 0 {
		from = 0
	}
	runes := []rune(c.text)
	if to > len(runes) {
		to = len(runes)
	}
	if from > len(runes) {
		from = len(runes)
	}
	return from, to
}

func (c *NumBox) clearSelection() {
	c.selectionAnchorPos = c.cursorPos
	c.selectionEndPos = c.cursorPos
}

func (c *NumBox) ProcessPropChange(key string, value interface{}) {
	if c.propIsProcessing {
		return
	}
	c.propIsProcessing = true
	defer func() { c.propIsProcessing = false }()

	switch key {
	case "decimals":
		c.SetDecimals(c.GetPropInt("decimals", c.decimals))
	case "min":
		c.min = c.GetPropFloat64("min", c.min)
		c.SetValue(c.value)
	case "max":
		c.max = c.GetPropFloat64("max", c.max)
		c.SetValue(c.value)
	case "step":
		c.step = c.GetPropFloat64("step", c.step)
	case "value":
		// Update internal state without setting the prop again.
		v := c.GetPropFloat64("value", c.value)
		if !math.IsNaN(v) {
			v = c.clamp(v)
			changed := !c.hasValue || v != c.value
			c.value = v
			c.hasValue = true
			c.refreshTextFromValue()
			if changed {
				c.fireValueChanged()
			}
			UpdateMainForm()
		}
	}
}

func (c *NumBox) buttonRect() (x, y, w, h int) {
	btnW := c.Height()
	if btnW < 18 {
		btnW = 18
	}
	return c.Width() - btnW, 0, btnW, c.Height()
}

func (c *NumBox) textRect() (x, y, w, h int) {
	btnX, _, _, _ := c.buttonRect()
	return c.padding, 0, btnX - c.padding, c.Height()
}

func (c *NumBox) draw(cnv *Canvas) {
	backColor := c.BackgroundColorWithAddElevation(-1)
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	if c.IsFocused() {
		backColor = c.BackgroundColorWithAddElevation(4)
	}

	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)

	// Spin buttons area
	btnX, btnY, btnW, btnH := c.buttonRect()
	btnBack := c.BackgroundColorWithAddElevation(1)
	cnv.FillRect(btnX, btnY, btnW, btnH, btnBack)
	cnv.FillRect(btnX, btnY+btnH/2, btnW, 1, c.BackgroundColorWithAddElevation(8))

	// Arrows
	arrowCol := c.ForegroundColor()
	// Draw simple triangles with lines so it works with any font.
	{
		cx := btnX + btnW/2
		midY := btnY + btnH/4
		size := btnW / 7
		if btnH/2/7 < size {
			size = btnH / 14
		}
		if size < 3 {
			size = 3
		}
		// Up triangle
		cnv.DrawLine(cx, midY-size, cx-size, midY+size, 1, arrowCol)
		cnv.DrawLine(cx, midY-size, cx+size, midY+size, 1, arrowCol)
		cnv.DrawLine(cx-size, midY+size, cx+size, midY+size, 1, arrowCol)
	}
	{
		cx := btnX + btnW/2
		midY := btnY + (btnH*3)/4
		size := btnW / 7
		if btnH/2/7 < size {
			size = btnH / 14
		}
		if size < 3 {
			size = 3
		}
		// Down triangle
		cnv.DrawLine(cx, midY+size, cx-size, midY-size, 1, arrowCol)
		cnv.DrawLine(cx, midY+size, cx+size, midY-size, 1, arrowCol)
		cnv.DrawLine(cx-size, midY-size, cx+size, midY-size, 1, arrowCol)
	}

	// Text
	tx, ty, tw, th := c.textRect()

	// Selection highlight (single line)
	if c.IsFocused() && c.hasSelection() {
		from, to := c.selectionRange()
		runes := []rune(c.text)
		charPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), string(runes))
		if err == nil && len(charPos) > 0 && from >= 0 && to >= 0 && from < len(charPos) && to < len(charPos) {
			x1 := tx + charPos[from]
			x2 := tx + charPos[to]
			if x2 < x1 {
				x1, x2 = x2, x1
			}
			_, lineH, _ := MeasureText(c.FontFamily(), c.FontSize(), "Q")
			selY := (c.Height() - lineH) / 2
			cnv.FillRect(x1, selY, x2-x1, lineH, c.BackgroundColorWithAddElevation(10))
		}
	}

	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.ForegroundColor())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(tx, ty, tw, th, c.text)

	// Cursor
	if c.IsFocused() && !c.hasSelection() {
		runes := []rune(c.text)
		if c.cursorPos < 0 {
			c.cursorPos = 0
		}
		if c.cursorPos > len(runes) {
			c.cursorPos = len(runes)
		}
		charPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), string(runes))
		if err == nil && len(charPos) > 0 {
			curX := tx + charPos[c.cursorPos]
			_, lineH, _ := MeasureText(c.FontFamily(), c.FontSize(), "Q")
			curY := (c.Height() - lineH) / 2
			cnv.FillRect(curX, curY, 1, lineH, c.ForegroundColor())
		}
	}
}

func (c *NumBox) onMouseDownLeft(x, y int, mods nuikey.KeyModifiers) {
	btnX, _, btnW, btnH := c.buttonRect()
	if x >= btnX && x < btnX+btnW {
		// Button click
		c.mouseButtonPressed = false
		c.draggingSelection = false
		c.clearSelection()
		if y < btnH/2 {
			c.stepBy(+c.Step())
		} else {
			c.stepBy(-c.Step())
		}
		return
	}

	// Place cursor near click position in text area.
	tx, _, _, _ := c.textRect()
	clickX := x - tx
	if clickX < 0 {
		clickX = 0
	}
	pos := c.cursorPosFromPixel(clickX)
	c.cursorPos = pos
	c.mouseButtonPressed = true
	c.draggingSelection = true
	if mods.Shift {
		c.selectionEndPos = c.cursorPos
	} else {
		c.selectionAnchorPos = c.cursorPos
		c.selectionEndPos = c.cursorPos
	}
	UpdateMainForm()
}

func (c *NumBox) onMouseMove(x, y int, mods nuikey.KeyModifiers) {
	_ = y
	_ = mods

	// Cursor shape depends on area.
	btnX, _, btnW, _ := c.buttonRect()
	if x >= btnX && x < btnX+btnW {
		c.SetMouseCursor(nuimouse.MouseCursorPointer)
	} else {
		c.SetMouseCursor(nuimouse.MouseCursorIBeam)
	}

	if !c.mouseButtonPressed || !c.draggingSelection {
		return
	}
	tx, _, _, _ := c.textRect()
	clickX := x - tx
	if clickX < 0 {
		clickX = 0
	}
	pos := c.cursorPosFromPixel(clickX)
	c.cursorPos = pos
	c.selectionEndPos = pos
	UpdateMainForm()
}

func (c *NumBox) onMouseUpLeft(x, y int, mods nuikey.KeyModifiers) {
	_ = x
	_ = y
	_ = mods
	c.mouseButtonPressed = false
	c.draggingSelection = false
	UpdateMainForm()
}

func (c *NumBox) onMouseWheel(deltaY int) {
	if deltaY == 0 {
		return
	}
	// When user is editing, first commit to keep wheel behavior predictable.
	c.commitText(true)

	step := c.Step()
	if step == 0 || math.IsNaN(step) || math.IsInf(step, 0) {
		return
	}
	delta := step * float64(deltaY)
	c.stepBy(delta)
}

func (c *NumBox) cursorPosFromPixel(px int) int {
	runes := []rune(c.text)
	charPos, err := GetCharPositions(c.FontFamily(), c.FontSize(), string(runes))
	if err != nil || len(charPos) == 0 {
		return len(runes)
	}
	if len(charPos) == 1 {
		return 0
	}
	if px <= (charPos[1]-charPos[0])/2 {
		return 0
	}
	for i := 1; i < len(charPos)-1; i++ {
		left := charPos[i] - (charPos[i]-charPos[i-1])/2
		right := charPos[i] + (charPos[i+1]-charPos[i])/2
		if px >= left && px < right {
			return i
		}
	}
	lastW := 0
	if len(charPos) > 1 {
		lastW = charPos[len(charPos)-1] - charPos[len(charPos)-2]
	}
	if px >= charPos[len(charPos)-1]-lastW/2 {
		return len(charPos) - 1
	}
	return len(runes)
}

func (c *NumBox) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	// basic selection operations
	if mods.Ctrl && key == nuikey.KeyA {
		c.selectionAnchorPos = 0
		c.selectionEndPos = len([]rune(c.text))
		c.cursorPos = c.selectionEndPos
		UpdateMainForm()
		return true
	}

	switch key {
	case nuikey.KeyArrowLeft:
		if c.cursorPos > 0 {
			c.cursorPos--
		}
		if mods.Shift {
			c.selectionEndPos = c.cursorPos
		} else {
			c.clearSelection()
		}
		UpdateMainForm()
		return true
	case nuikey.KeyArrowRight:
		if c.cursorPos < len([]rune(c.text)) {
			c.cursorPos++
		}
		if mods.Shift {
			c.selectionEndPos = c.cursorPos
		} else {
			c.clearSelection()
		}
		UpdateMainForm()
		return true
	case nuikey.KeyHome:
		c.cursorPos = 0
		if mods.Shift {
			c.selectionEndPos = c.cursorPos
		} else {
			c.clearSelection()
		}
		UpdateMainForm()
		return true
	case nuikey.KeyEnd:
		c.cursorPos = len([]rune(c.text))
		if mods.Shift {
			c.selectionEndPos = c.cursorPos
		} else {
			c.clearSelection()
		}
		UpdateMainForm()
		return true
	case nuikey.KeyBackspace:
		c.deleteLeft()
		return true
	case nuikey.KeyDelete:
		c.deleteRight()
		return true
	case nuikey.KeyArrowUp:
		c.stepBy(+c.Step())
		return true
	case nuikey.KeyArrowDown:
		c.stepBy(-c.Step())
		return true
	case nuikey.KeyEnter:
		c.commitText(true)
		return true
	}
	return false
}

func (c *NumBox) onChar(ch rune, mods nuikey.KeyModifiers) bool {
	_ = mods
	if ch < 32 {
		return false
	}

	// Allow digits, one decimal separator, and leading minus.
	if ch >= '0' && ch <= '9' {
		c.insertRune(ch)
		return true
	}

	if ch == '.' || ch == ',' {
		// normalize to '.'
		c.insertRune('.')
		return true
	}

	if ch == '-' {
		c.insertMinus()
		return true
	}

	return true
}

func (c *NumBox) insertRune(ch rune) {
	if c.hasSelection() {
		c.deleteSelection()
	}
	runes := []rune(c.text)
	if c.cursorPos < 0 {
		c.cursorPos = 0
	}
	if c.cursorPos > len(runes) {
		c.cursorPos = len(runes)
	}

	// Only one dot allowed.
	if ch == '.' {
		for _, r := range runes {
			if r == '.' {
				return
			}
		}
		// If decimals==0, still allow typing dot (user may change decimals later),
		// but it will be normalized on commit.
	}

	left := append([]rune{}, runes[:c.cursorPos]...)
	right := append([]rune{}, runes[c.cursorPos:]...)
	runes = append(left, ch)
	runes = append(runes, right...)

	c.text = string(runes)
	c.cursorPos++
	c.clearSelection()
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) insertMinus() {
	if c.hasSelection() {
		c.deleteSelection()
	}
	runes := []rune(c.text)
	if len(runes) > 0 && runes[0] == '-' {
		// toggle off
		runes = runes[1:]
		if c.cursorPos > 0 {
			c.cursorPos--
		}
	} else {
		runes = append([]rune{'-'}, runes...)
		if c.cursorPos >= 0 {
			c.cursorPos++
		}
	}
	c.text = string(runes)
	c.clearSelection()
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) deleteSelection() {
	if !c.hasSelection() {
		return
	}
	from, to := c.selectionRange()
	runes := []rune(c.text)
	if from == to {
		return
	}
	runes = append(runes[:from], runes[to:]...)
	c.text = string(runes)
	c.cursorPos = from
	c.clearSelection()
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) deleteLeft() {
	if c.hasSelection() {
		c.deleteSelection()
		return
	}
	runes := []rune(c.text)
	if c.cursorPos <= 0 || len(runes) == 0 {
		return
	}
	runes = append(runes[:c.cursorPos-1], runes[c.cursorPos:]...)
	c.cursorPos--
	c.text = string(runes)
	c.clearSelection()
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) deleteRight() {
	if c.hasSelection() {
		c.deleteSelection()
		return
	}
	runes := []rune(c.text)
	if c.cursorPos < 0 || c.cursorPos >= len(runes) || len(runes) == 0 {
		return
	}
	runes = append(runes[:c.cursorPos], runes[c.cursorPos+1:]...)
	c.text = string(runes)
	c.clearSelection()
	c.editingText = true
	UpdateMainForm()
}

func (c *NumBox) parseText(text string) (float64, bool) {
	s := strings.TrimSpace(text)
	if s == "" || s == "-" || s == "." || s == "-." {
		return 0, false
	}
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

func (c *NumBox) commitText(force bool) {
	if !force && !c.editingText {
		return
	}
	if v, ok := c.parseText(c.text); ok {
		c.SetValue(v)
		return
	}
	// If invalid, revert to current value formatting.
	c.refreshTextFromValue()
	UpdateMainForm()
}

func (c *NumBox) stepBy(delta float64) {
	if delta == 0 || math.IsNaN(delta) || math.IsInf(delta, 0) {
		return
	}
	// If user is editing, commit current text first so stepping is predictable.
	c.commitText(true)

	v := c.clamp(c.value + delta)
	// Round to decimals to stabilize drift.
	if c.decimals >= 0 && c.decimals <= 16 {
		p := math.Pow10(c.decimals)
		v = math.Round(v*p) / p
	}
	c.SetValue(v)
}

