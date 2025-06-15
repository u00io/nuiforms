package ex02widget

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func NewMyWidget(col color.RGBA) *ui.Widget {
	c := ui.NewWidget()
	c.SetSize(100, 100)
	c.SetOnPaint(func(cnv *ui.Canvas) {
		cnv.SetColor(col)
		cnv.FillRect(0, 0, c.Width(), c.Height(), col)
		cnv.DrawText(10, 40, "MyWidget", "robotomono", 16, ui.DefaultForeground, false)
	})

	btn1 := ui.NewButton()
	btn1.SetPosition(8, 8)
	c.AddWidget(btn1)

	return c
}
