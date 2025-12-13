package ex03masterdetail

import (
	"fmt"
	"image/color"

	_ "embed"

	"github.com/u00io/nuiforms/ui"
)

type Tab1Widget struct {
	ui.Widget

	btn1 *ui.Button
}

// embed design time metadata

//go:embed design.xml
var design string

func NewTab1Widget() *Tab1Widget {
	var c Tab1Widget
	c.InitWidget()
	c.SetBackgroundColor(color.RGBA{R: 150, G: 50, B: 50, A: 255})
	c.btn1 = ui.NewButton("button1")
	fmt.Println(design)
	return &c
}
