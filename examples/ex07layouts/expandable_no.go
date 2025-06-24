package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

// Example: Button

type ExpandableNo struct {
	ui.Widget
}

func NewExpandableNo() *ExpandableNo {
	var c ExpandableNo
	c.InitWidget()
	c.SetXExpandable(false)
	c.SetYExpandable(false)
	c.SetMaxWidth(100)
	c.SetMaxHeight(30)
	c.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableNo) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Width(), c.Height(), color.RGBA{R: 255, G: 255, B: 0, A: 255})
}
