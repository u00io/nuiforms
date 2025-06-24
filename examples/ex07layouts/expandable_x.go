package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

// example: TextBox
type ExpandableX struct {
	ui.Widget
}

func NewExpandableX() *ExpandableX {
	var c ExpandableX
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(false)
	c.SetMaxHeight(30)
	c.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableX) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Width(), c.Height(), color.RGBA{R: 255, G: 0, B: 0, A: 255})
}
