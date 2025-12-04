package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Button struct {
	Widget

	fontSize float64

	pressed       bool
	text          string
	onButtonClick func(btn *Button)
}

func NewButton(text string) *Button {
	var c Button
	c.InitWidget()
	c.SetTypeName("Button")
	c.SetMinSize(100, 30)
	c.SetMaxSize(10000, 30)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("Button")
	c.SetCanBeFocused(true)

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)
	c.SetOnKeyDown(c.onKeyDown)

	c.SetText(text)

	c.fontSize = -1

	return &c
}

func (c *Button) buttonFontSize() float64 {
	if c.fontSize > 0 {
		return c.fontSize
	}
	return c.FontSize()
}

func (c *Button) SetFontSize(size float64) {
	c.fontSize = size
	UpdateMainForm()
}

func (c *Button) Text() string {
	return c.text
}

func (c *Button) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *Button) SetOnButtonClick(fn func(btn *Button)) {
	c.onButtonClick = fn
}

func (c *Button) Press() {
	if c.onButtonClick != nil {
		c.onButtonClick(c)
	}
}

func (c *Button) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter || key == nuikey.KeySpace {
		c.Press()
	}
	return true
}

func (c *Button) draw(cnv *Canvas) {
	backColor := c.BackgroundColor()
	if c.IsHovered() && c.enabled {
		backColor = c.BackgroundColorAccent1()
	}
	if c.pressed {
		backColor = c.BackgroundColorAccent2()
	}
	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)

	cnv.SetHAlign(HAlignCenter)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.Color())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.buttonFontSize())
	cnv.DrawText(0, 0, c.Width(), c.Height(), c.text)

	cnv.SetColor(c.BackgroundColorAccent2())
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}

func (c *Button) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.enabled == false {
		return false
	}
	c.pressed = true
	return true
}

func (c *Button) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.enabled == false {
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
		if c.onButtonClick != nil {
			c.onButtonClick(c)
		}
	}

	return true
}
