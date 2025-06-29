package ex06font

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type FontExampleWidget struct {
	ui.Widget
}

func NewFontExampleWidget() *FontExampleWidget {
	var widget FontExampleWidget
	widget.InitWidget()
	widget.SetTypeName("FontExampleWidget")
	widget.SetOnPaint(widget.draw)
	return &widget
}

func (c *FontExampleWidget) draw(cnv *ui.Canvas) {
	cnv.DrawLine(100, 100, 100, 150, 1, color.RGBA{R: 0, G: 200, B: 0, A: 255})
	cnv.DrawLine(100, 100, 200, 100, 1, color.RGBA{R: 0, G: 200, B: 0, A: 255})
}

func Run(form *ui.Form) {
	widget := NewFontExampleWidget()
	form.Panel().AddWidgetOnGrid(widget, 0, 0)
}
