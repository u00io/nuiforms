package ex03masterdetail

import (
	"github.com/u00io/nuiforms/ui"
)

type Tab2Widget struct {
	ui.Widget
}

func NewTab2Widget(name string) *Tab2Widget {
	var c Tab2Widget
	c.InitWidget()
	//c.widget.SetBackgroundColor(color.RGBA{R: 50, G: 150, B: 50, A: 255})
	table := ui.NewTable()
	table.SetColumnCount(3)
	table.SetRowCount(10)
	c.AddWidgetOnGrid(table, 0, 0)

	table.SetCellText2(1, 1, "Column 1")
	return &c
}
