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

////////////////////////////////////////////////////////////////////////////////////
// NewLabel creates a new Label widget with the specified text.

func NewLabel(text string) *Label {
	var c Label
	c.InitWidget()
	c.SetTypeName("Label")
	c.SetMaxWidth(500)
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) {
		fmt.Println("Label clicked:", c.text)
	})
	c.SetOnPaint(func(cnv *Canvas) {
		cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignLeft, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	c.SetText(text)
	return &c
}

/////////////////////////////////////////////////////////////////////////////////
// Label methods

func (c *Label) Text() string {
	return c.text
}

func (c *Label) SetText(text string) {
	c.text = text
	c.updateInnerSize()
	UpdateMainForm()
}

/////////////////////////////////////////////////////////////////////////////////
// Label private methods

func (c *Label) updateInnerSize() {
	_, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), "0")
	if err != nil {
		return
	}
	c.innerHeight = textHeight * 1

	var maxTextWidth int
	textWidth, _, err := MeasureText(c.FontFamily(), c.FontSize(), c.text)
	if err != nil {
		return
	}
	if textWidth > maxTextWidth {
		maxTextWidth = textWidth
	}
	c.innerWidth = maxTextWidth
	//c.innerHeight = 0
	c.SetMinSize(c.innerWidth, c.innerHeight)
}
