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
		cnv.SetVAlign(ui.VAlignCenter)
		cnv.SetHAlign(ui.HAlignCenter)
		cnv.SetColor(c.ForegroundColor())
		cnv.SetFontFamily(c.FontFamily())
		cnv.SetFontSize(c.FontSize())
		cnv.DrawText(0, 0, c.Width(), c.Height(), name)
	})

	return &c
}
