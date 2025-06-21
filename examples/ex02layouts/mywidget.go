package ex02layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type MyWidget struct {
	widget ui.Widget
}

func NewMyWidget(name string) *MyWidget {
	var c MyWidget
	c.widget.InitWidget()
	c.widget.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})

	c.widget.SetOnPaint(func(cnv *ui.Canvas) {
		cnv.DrawTextMultiline(0, 0, c.widget.Width(), c.widget.Height(), ui.HAlignCenter, ui.VAlignCenter, name, color.White, "robotomono", 18, false)
	})

	return &c
}

func (c *MyWidget) Widgeter() any {
	return &c.widget
}
