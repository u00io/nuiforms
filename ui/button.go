package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Button struct {
	Widget
	pressed bool
}

const (
	DefaultButtonMinWidth = 100
)

func NewButton(text string) *Button {
	var c Button
	c.InitWidget()

	c.dontAllowOnClickIfDisabled = true

	c.SetTypeName("Button")
	c.SetMinSize(DefaultButtonMinWidth, DefaultUiLineHeight)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("Button")
	c.SetCanBeFocused(true)
	c.SetElevation(3)

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)
	// c.SetOnKeyDown(c.onKeyDown)

	c.SetText(text)
	c.SetProp("padding", 6)
	// SetProp only dispatches ProcessPropChange (which would normally grow
	// the button to fit this text) once the button is attached to a form -
	// c.form is still nil here, during construction, so that dispatch is a
	// no-op. Size from the text directly so the very first, most natural
	// construction path (NewButton(text) then AddWidget(...)) is sized
	// correctly too, not just a later SetText call made after attaching.
	c.updateSizeFromText()

	return &c
}

func (c *Button) Text() string {
	return c.GetPropString("text", "")
}

func (c *Button) SetText(text string) {
	c.SetProp("text", text)
	c.form.Update()
}

// Push simulates a click on the button, triggering the "onclick" handler
// programmatically (used e.g. by Form's accept/cancel buttons on Enter/Esc).
func (c *Button) Push() {
	if !c.Enabled() {
		return
	}
	if f := c.GetPropFunction("onclick"); f != nil {
		f()
	}
}

func (c *Button) ProcessKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	return c.onKeyDown(key, mods)
}

func (c *Button) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter || key == nuikey.KeySpace {
		if f := c.GetPropFunction("onclick"); f != nil {
			f()
			return true
		}
	}
	return false
}

func (c *Button) draw(cnv *Canvas) {
	backColor := c.BackgroundColor()

	if c.IsHovered() && c.Enabled() {
		backColor = c.BackgroundColorWithAddElevation(1)
	}
	if c.pressed {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	_ = backColor

	cnv.SetColor(backColor)
	cnv.FillRoundedRect(0, 0, c.Width(), c.Height(), 5)

	foreColor := c.ForegroundColor()
	if !c.Enabled() {
		foreColor = c.ForegroundColorDisabled()
	}

	cnv.SetHAlign(HAlignCenter)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(foreColor)
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(0, 0, c.Width(), c.Height(), c.Text())

	if c.IsFocused() {
		cnv.SetColor(ColorFromHex("#0088AA"))
		cnv.DrawRoundedRect(1, 1, c.Width()-2, c.Height()-2, 4)
	}
}

func (c *Button) ProcessPropChange(key string, value interface{}) {
	if key == "enabled" {
		if c.Enabled() {
			c.SetMouseCursor(nuimouse.MouseCursorPointer)
		} else {
			c.SetMouseCursor(nuimouse.MouseCursorNotDefined)
		}
	}

	c.updateSizeFromText()
	c.form.UpdateLayout()
}

// updateSizeFromText grows the button's min size to fit its current text.
// Compares against c.minWidth directly rather than the cached MinWidth()
// getter, since SetMinWidth doesn't invalidate that cache - reading it here
// (e.g. on the very first call, from NewButton before any layout pass has
// run) would otherwise lock in a stale value once a real layout pass reads
// MinWidth() later.
func (c *Button) updateSizeFromText() {
	padding := c.GetPropInt("padding", 6)

	textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), c.Text())
	if err != nil {
		return
	}
	if textHeight < DefaultUiLineHeight {
		textHeight = DefaultUiLineHeight
	}

	c.SetMinHeight(textHeight)

	if minWidth := textWidth + padding*2; minWidth > c.minWidth {
		c.SetMinWidth(minWidth)
	}
}

func (c *Button) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if !c.Enabled() {
		return false
	}
	c.pressed = true
	return true
}

func (c *Button) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if !c.Enabled() {
		return false
	}

	wasPressed := c.pressed
	c.pressed = false

	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// Released outside the button - cancel the click.
		return false
	}

	if wasPressed {
		if f := c.GetPropFunction("onclick"); f != nil {
			f()
		}
	}

	return true
}
