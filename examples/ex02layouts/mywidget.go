package ex02layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type MyWidget struct {
	ui.Widget
}

func NewMyWidget(name string) *MyWidget {
	var c MyWidget
	c.InitWidget()
	c.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})

	c.SetOnPaint(func(cnv *ui.Canvas) {
		cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), ui.HAlignCenter, ui.VAlignCenter, name, color.White, "robotomono", 18, false)
	})

	return &c
}
