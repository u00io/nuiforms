package ex03masterdetail

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func NewTab2Widget() *ui.Widget {
	col := color.RGBA{R: 0x00, G: 0xCC, B: 0xFF, A: 0xFF} // Blue color
	c := ui.NewWidget()
	c.SetSize(100, 100)
	c.SetBackgroundColor(col)

	btn1 := ui.NewButton()
	btn1.SetPosition(8, 8)
	c.AddWidget(btn1)

	return c
}
