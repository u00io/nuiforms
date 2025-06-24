package ex07layouts

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type ExpandableY struct {
	ui.Widget
}

func NewExpandableY() *ExpandableY {
	var c ExpandableY
	c.InitWidget()
	c.SetXExpandable(false)
	c.SetYExpandable(true)
	c.SetMaxWidth(100)
	c.SetOnPaint(c.draw)
	c.AddTimer(1000, func() {
		fmt.Println("ExpandableY timer tick", c.Width(), c.Height())
	})
	return &c
}

func (c *ExpandableY) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Width(), c.Height(), color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
