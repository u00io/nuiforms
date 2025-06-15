package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

// example: TextBox
type ExpandableX struct {
	Widget ui.Widget
}

func NewExpandableX() *ExpandableX {
	var c ExpandableX
	c.Widget.InitWidget()
	c.Widget.SetXExpandable(true)
	c.Widget.SetYExpandable(false)
	c.Widget.SetMaxHeight(30)
	c.Widget.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableX) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Widget.Width(), c.Widget.Height(), color.RGBA{R: 255, G: 0, B: 0, A: 255})
}
