package ex12cards

import (
	"image/color"

	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	mainWidget := NewCards()
	form.Panel().AddWidgetOnGrid(mainWidget, 0, 0)
}

type Card struct {
	ui.Widget
	categoryName string
}

func NewCategoryWidget(categoryName string) *Card {
	var c Card
	c.InitWidget()
	c.SetTypeName("Card")
	c.categoryName = categoryName
	c.AddWidgetOnGrid(ui.NewLabel(categoryName), 0, 0)
	c.SetYExpandable(false)
	c.SetMinHeight(120)
	c.SetMaxHeight(120)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	return &c
}
