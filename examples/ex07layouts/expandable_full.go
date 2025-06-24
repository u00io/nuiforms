package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type ExpandableFull struct {
	ui.Widget
	col color.RGBA
}

func NewExpandableFull(col color.RGBA) *ExpandableFull {
	var c ExpandableFull
	c.InitWidget()
	c.col = col
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	c.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableFull) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.Width(), c.Height(), c.col)
}

func (c *ExpandableFull) AddWidgetOnGrid(w ui.Widgeter, row, col int) {
	c.AddWidgetOnGrid(w, row, col)
}
