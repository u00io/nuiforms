package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type tabWidgetButton struct {
	Widget

	pressed       bool
	text          string
	onButtonClick func(btn *tabWidgetButton)
}

func NewTabWidgetButton(text string) *tabWidgetButton {
	var c tabWidgetButton
	c.InitWidget()
	c.SetTypeName("TabWidgetButton")
	c.SetMinSize(100, 30)
	c.SetMaxSize(10000, 30)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("Button")

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)

	c.SetText(text)

	c.SetCanBeFocused(true)

	return &c
}

func (c *tabWidgetButton) Text() string {
	return c.text
}

func (c *tabWidgetButton) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *tabWidgetButton) SetOnButtonClick(fn func(btn *tabWidgetButton)) {
	c.onButtonClick = fn
}

func (c *tabWidgetButton) draw(cnv *Canvas) {
	cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
}

func (c *tabWidgetButton) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = true
	return true
}

func (c *tabWidgetButton) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
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
