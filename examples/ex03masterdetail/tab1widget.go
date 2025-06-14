package ex03masterdetail

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func NewTab1Widget() *ui.Widget {
	col := color.RGBA{R: 0xFF, G: 0xCC, B: 0x00, A: 0xFF} // Yellow color
	c := ui.NewWidget()
	c.SetSize(100, 100)
	c.SetBackgroundColor(col)

	btn1 := ui.NewButton()
	btn1.SetSize(80, 30)
	btn1.SetPosition(8, 8)
	btn1.SetProp("text", "OK")
	btn1.SetProp("onClick", func() {
		fmt.Println("Button OK clicked")
	})
	c.AddWidget(btn1)

	return c
}
