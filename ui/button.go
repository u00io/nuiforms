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
	c.SetOnKeyDown(c.onKeyDown)

	c.SetText(text)
	c.SetProp("padding", 6)

	return &c
}

func (c *Button) Text() string {
	return c.GetPropString("text", "")
}

func (c *Button) SetText(text string) {
	c.SetProp("text", text)
	UpdateMainForm()
}

func (c *Button) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter || key == nuikey.KeySpace {
		if f := c.GetPropFunction("onclick"); f != nil {
			f()
		}
	}
	return true
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
}

func (c *Button) ProcessPropChange(key string, value interface{}) {
	padding := c.GetPropInt("padding", 6)
	if key == "enabled" {
		if c.Enabled() {
			c.SetMouseCursor(nuimouse.MouseCursorPointer)
		} else {
			c.SetMouseCursor(nuimouse.MouseCursorNotDefined)
		}
	}

	textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), c.Text())
	if err != nil {
		return
	}
	if textHeight < DefaultUiLineHeight {
		textHeight = DefaultUiLineHeight
	}
	_ = textHeight
	_ = textWidth
	_ = padding

	c.SetMinHeight(textHeight)

	if textWidth+padding*2 > c.MinWidth() {
		c.SetMinWidth(textWidth + padding*2)
	}

	UpdateMainFormLayout()
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
	c.pressed = false

	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// MouseUp outside the button area, ignore
		return false
	}

	hoverWidgeter := MainForm.hoverWidget
	var localWidgeter Widgeter = c
	if hoverWidgeter == localWidgeter {
		//c.Press()
	}

	return false
}
