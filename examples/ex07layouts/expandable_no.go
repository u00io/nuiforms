package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

// Example: Button

type ExpandableNo struct {
	Widget ui.Widget
}

func NewExpandableNo() *ExpandableNo {
	var c ExpandableNo
	c.Widget.InitWidget()
	c.Widget.SetXExpandable(false)
	c.Widget.SetYExpandable(false)
	c.Widget.SetMaxWidth(100)
	c.Widget.SetMaxHeight(30)
	c.Widget.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableNo) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Widget.Width(), c.Widget.Height(), color.RGBA{R: 255, G: 255, B: 0, A: 255})
}
