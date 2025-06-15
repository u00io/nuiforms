package ex03masterdetail

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func NewMasterWidget() *ui.Widget {
	c := ui.NewWidget()

	c.SetBackgroundColor(color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xFF}) // Light gray background

	{
		btnOpenTab1 := ui.NewButton()
		btnOpenTab1.SetPosition(0, 0)
		btnOpenTab1.SetAnchors(true, false, true, false)
		c.AddWidget(btnOpenTab1)
	}

	{
		btnOpenTab2 := ui.NewButton()
		btnOpenTab2.SetPosition(0, 30)
		btnOpenTab2.SetAnchors(true, false, true, false)
		c.AddWidget(btnOpenTab2)
	}

	return c
}
