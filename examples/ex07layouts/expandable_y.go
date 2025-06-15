package ex07layouts

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type ExpandableY struct {
	Widget ui.Widget
}

func NewExpandableY() *ExpandableY {
	var c ExpandableY
	c.Widget.InitWidget()
	c.Widget.SetXExpandable(false)
	c.Widget.SetYExpandable(true)
	c.Widget.SetMaxWidth(100)
	c.Widget.SetOnPaint(c.draw)
	c.Widget.AddTimer(1000, func() {
		fmt.Println("ExpandableY timer tick", c.Widget.Width(), c.Widget.Height())
	})
	return &c
}

func (c *ExpandableY) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Widget.Width(), c.Widget.Height(), color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
