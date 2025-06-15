package ui

import (
	"fmt"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Label struct {
	Widget
	text string
}

func NewLabel(text string) *Label {
	var c Label
	c.text = text
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) {
		fmt.Println("Label clicked:", c.text)
	})
	c.SetOnPaint(func(cnv *Canvas) {
		cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignCenter, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	return &c
}

func (c *Label) SetText(text string) {
	c.text = text
	UpdateMainForm()
}
