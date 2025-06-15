package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type ExpandableFull struct {
	widget ui.Widget
	col    color.RGBA
}

func (c *ExpandableFull) Widgeter() any {
	return &c.widget
}

func NewExpandableFull(col color.RGBA) *ExpandableFull {
	var c ExpandableFull
	c.widget.InitWidget()
	c.col = col
	c.widget.SetXExpandable(true)
	c.widget.SetYExpandable(true)
	c.widget.SetOnPaint(c.draw)
	return &c
}

func (c *ExpandableFull) draw(cnv *ui.Canvas) {
	cnv.FillRect(0, 0, c.widget.Width(), c.widget.Height(), c.col)
}

func (c *ExpandableFull) AddWidgetOnGrid(w any, row, col int) {
	c.widget.AddWidgetOnGrid(w, row, col)
}
