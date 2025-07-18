package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Button struct {
	Widget

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

	c.SetText(text)

	return &c
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

func (c *Button) draw(cnv *Canvas) {
	backColor := GetThemeColor("background", DefaultBackground)
	if c.IsHovered() {
		backColor = GetThemeColor("button.background.hover", DefaultBackground)
	}

	if c.pressed {
		backColor = GetThemeColor("button.background.pressed", DefaultBackground)
	}

	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)
	cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
}

func (c *Button) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = true
	return true
}

func (c *Button) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
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
