package ex03masterdetail

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type Tab1Widget struct {
	ui.Widget
}

func NewTab1Widget() *Tab1Widget {
	var c Tab1Widget
	c.InitWidget()
	c.SetBackgroundColor(color.RGBA{R: 150, G: 50, B: 50, A: 255})
	return &c
}
