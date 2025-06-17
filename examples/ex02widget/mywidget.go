package ex02widget

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type MyWidget struct {
	widget ui.Widget
}

func NewMyWidget(col color.RGBA) *MyWidget {
	var c MyWidget
	c.widget.InitWidget()

	c.widget.SetOnPaint(func(cnv *ui.Canvas) {
		cnv.SetColor(col)
		cnv.FillRect(0, 0, c.widget.Width()+50, c.widget.Height()+50, col)
		cnv.DrawTextMultiline(0, 0, 100, 100, ui.HAlignLeft, ui.VAlignTop, "1234567890", color.White, "robotomono", 18, false)
	})

	return &c
}

func (c *MyWidget) Widgeter() any {
	return &c.widget
}
