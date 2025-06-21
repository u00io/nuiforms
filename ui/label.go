package ui

import (
	"fmt"

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
	c.text = text
	c.widget.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) {
		fmt.Println("Label clicked:", c.text)
	})
	c.widget.SetOnPaint(func(cnv *Canvas) {
		cnv.DrawTextMultiline(0, 0, c.widget.Width(), c.widget.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	return &c
}

func (c *Label) SetText(text string) {
	c.text = text
	UpdateMainForm()
}
