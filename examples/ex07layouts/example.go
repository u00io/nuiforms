package ex07layouts

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	w00 := NewExpandableFull(color.RGBA{R: 10, G: 10, B: 10, A: 255})
	w01 := NewExpandableFull(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	w10 := NewExpandableFull(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	w11 := NewExpandableFull(color.RGBA{R: 150, G: 150, B: 150, A: 255})

	{
		expX := NewExpandableX()
		w00.AddWidgetOnGrid(expX, 0, 0)
	}

	{
		expNo := NewExpandableNo()
		w01.AddWidgetOnGrid(expNo, 0, 0)
	}

	{
		expY := NewExpandableY()
		w10.AddWidgetOnGrid(expY, 0, 0)
	}

	{
		expFull := NewExpandableFull(color.RGBA{R: 200, G: 200, B: 200, A: 255})
		w11.AddWidgetOnGrid(expFull, 0, 0)
	}

	form.Panel().AddWidgetOnGrid(w00, 0, 0)
	form.Panel().AddWidgetOnGrid(w01, 0, 1)
	form.Panel().AddWidgetOnGrid(w10, 1, 0)
	form.Panel().AddWidgetOnGrid(w11, 1, 1)
	form.Exec()
}
