package ex13table

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	mainWidget := NewTableWidget()
	form.Panel().AddWidgetOnGrid(mainWidget, 0, 0)
}

type TableWidget struct {
	ui.Widget
	lvItem *ui.Table

	lblInnerWidget    *ui.Label
	txtBoxInnerWidget *ui.TextBox
}

func NewTableWidget() *TableWidget {
	var c TableWidget
	c.InitWidget()
	c.SetTypeName("TableWidget")
	c.lvItem = ui.NewTable()
	c.lvItem.SetColumnCount(3)
	c.lvItem.SetColumnWidth(0, 200)
	c.lvItem.SetColumnWidth(1, 200)
	c.lvItem.SetColumnWidth(2, 200)
	c.lvItem.SetColumnName(0, "Col1")
	c.lvItem.SetColumnName(1, "Col2")
	c.lvItem.SetColumnName(2, "Col3")
	c.lvItem.SetRowCount(30)
	for i := 0; i < 30; i++ {
		c.lvItem.SetCellText(0, i, "row "+fmt.Sprint(i))
		c.lvItem.SetCellText(1, i, "col2 text")
		c.lvItem.SetCellText(2, i, "col2 text")
	}

	c.AddWidgetOnGrid(c.lvItem, 0, 0)
	c.SetYExpandable(false)

	c.lblInnerWidget = ui.NewLabel("This is an inner widget")
	c.lblInnerWidget.SetBackgroundColor(color.RGBA{R: 90, G: 90, B: 90, A: 255})
	c.lvItem.AddWidgetOnTable(c.lblInnerWidget, 1, 1, 2, 2)

	c.txtBoxInnerWidget = ui.NewTextBox()
	c.txtBoxInnerWidget.SetText("This is an inner TextBox")
	c.lvItem.AddWidgetOnTable(c.txtBoxInnerWidget, 2, 5, 1, 1)

	return &c
}
