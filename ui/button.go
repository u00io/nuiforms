package ui

import (
	"time"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Button struct {
	widget Widget

	pressed       bool
	text          string
	onButtonClick func(btn *Button)
}

func NewButton() *Button {
	var c Button
	c.widget.InitWidget()

	c.widget.SetName("Button-" + time.Now().Format("20060102150405"))
	c.widget.SetMinSize(100, 30)
	c.widget.SetMaxSize(100, 30)
	c.widget.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("Button")

	c.widget.SetOnPaint(c.draw)
	c.widget.SetOnMouseDown(c.buttonProcessMouseDown)
	c.widget.SetOnMouseUp(c.buttonProcessMouseUp)

	return &c
}

func (c *Button) Widgeter() any {
	return &c.widget
}

func (c *Button) SetPosition(x, y int) {
	c.widget.SetPosition(x, y)
}

func (c *Button) SetMaxSize(width, height int) {
	c.widget.SetMaxSize(width, height)
}

func (c *Button) SetAnchors(left, right, top, bottom bool) {
	c.widget.SetAnchors(left, right, top, bottom)
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
	if c.widget.IsHovered() {
		backColor = GetThemeColor("button.background.hover", DefaultBackground)
		if c.pressed {
			backColor = GetThemeColor("button.background.pressed", DefaultBackground)
		}
	}
	cnv.FillRect(0, 0, c.widget.Width(), c.widget.Height(), backColor)
	cnv.DrawTextMultiline(0, 0, c.widget.Width(), c.widget.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
}

func (c *Button) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	c.pressed = true
}

func (c *Button) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	c.pressed = false
	if MainForm.hoverWidget == &c.widget {
		if c.onButtonClick != nil {
			c.onButtonClick(c)
		}
	}
}
