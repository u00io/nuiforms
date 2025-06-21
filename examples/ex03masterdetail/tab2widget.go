package ex03masterdetail

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type Tab2Widget struct {
	widget ui.Widget
}

func (c *Tab2Widget) Widgeter() any {
	return &c.widget
}

func NewTab2Widget(name string) *Tab2Widget {
	var c Tab2Widget
	c.widget.InitWidget()
	c.widget.SetBackgroundColor(color.RGBA{R: 50, G: 150, B: 50, A: 255})
	lbl := ui.NewLabel(name)
	c.widget.AddWidgetOnGrid(lbl, 0, 0)
	return &c
}
