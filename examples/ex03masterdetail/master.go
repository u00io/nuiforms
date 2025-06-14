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
		btnOpenTab1.SetSize(c.W(), 30)
		btnOpenTab1.SetAnchors(true, false, true, false)
		btnOpenTab1.SetProp("text", "Tab 1")
		btnOpenTab1.SetProp("onClick", func() {
			f, ok := c.GetProp("open1").(func())
			if ok {
				f()
			}
		})
		c.AddWidget(btnOpenTab1)
	}

	{
		btnOpenTab2 := ui.NewButton()
		btnOpenTab2.SetPosition(0, 30)
		btnOpenTab2.SetSize(c.W(), 30)
		btnOpenTab2.SetAnchors(true, false, true, false)
		btnOpenTab2.SetProp("text", "Tab 2")
		btnOpenTab2.SetProp("onClick", func() {
			f, ok := c.GetProp("open2").(func())
			if ok {
				f()
			}
		})
		c.AddWidget(btnOpenTab2)
	}

	return c
}
