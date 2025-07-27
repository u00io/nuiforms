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
	c.lvItem.SetHeaderRowCount(2)

	c.lvItem.SetEditTriggerDoubleClick(true)
	c.lvItem.SetEditTriggerEnter(true)
	c.lvItem.SetEditTriggerF2(true)

	c.lvItem.SetColumnCellName(0, 1, "Col1 Header")
	c.lvItem.SetRowCount(10)
	for i := 0; i < 10; i++ {
		c.lvItem.SetCellText(0, i, "row "+fmt.Sprint(i))
		c.lvItem.SetCellText(1, i, "col2 text")
		c.lvItem.SetCellText(2, i, "col2 text")
	}

	c.lvItem.SetHeaderCellSpan(0, 0, 1, 2)
	c.lvItem.SetHeaderCellSpan(1, 0, 2, 1)

	c.lvItem.SetColumnCellName(1, 0, "SPANNED COLUMN")
	c.lvItem.SetColumnCellName(1, 1, "COL1 HEADER")
	c.lvItem.SetColumnCellName(2, 1, "COL2 HEADER")

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

type ValueSelector struct {
	ui.Widget
	items []*ValueSelectorItem
}

type ValueSelectorItem struct {
	text  string
	value interface{}
}

func NewValueSelector() *ValueSelector {
	var c ValueSelector
	c.InitWidget()
	return &c
}

func (c *ValueSelector) AddItem(text string, value interface{}) {
	var item ValueSelectorItem
	item.text = text
	item.value = value
	c.items = append(c.items, &item)
}

func (c *ValueSelector) Show(xAtForm int, yAtForm int) {
	c.SetPosition(xAtForm, yAtForm)
	ui.MainForm.Panel().AppendPopupWidget(c)
	ui.UpdateMainForm()
}
