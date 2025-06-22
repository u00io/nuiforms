package ui

import (
	"fmt"
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Label struct {
	widget Widget
	text   string
}

func (c *Label) Widgeter() any {
	return &c.widget
}

func NewLabel(text string) *Label {
	var c Label
	c.widget.InitWidget()
	c.widget.SetBackgroundColor(color.RGBA{0, 150, 200, 255})
	c.text = text
	c.widget.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) {
		fmt.Println("Label clicked:", c.text)
	})
	c.widget.SetOnPaint(func(cnv *Canvas) {
		cnv.DrawTextMultiline(0, 0, c.widget.Width(), c.widget.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	return &c
}

func (c *Label) SetName(name string) {
	c.widget.SetName(name)
}

func (c *Label) SetMinSize(width, height int) {
	c.widget.SetMinSize(width, height)
}

func (c *Label) SetMaxSize(width, height int) {
	c.widget.SetMaxSize(width, height)
}

func (c *Label) SetText(text string) {
	c.text = text
	UpdateMainForm()
}
