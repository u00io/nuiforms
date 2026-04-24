package ex13table

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	mainWidget := NewTableWidget()
	form.Panel().AddWidgetOnGrid(mainWidget, 0, 0)
}

//go:embed testImage.png
var testImage []byte

func TestImage() image.Image {
	im, _ := png.Decode(strings.NewReader(string(testImage)))
	return im
}

type TableWidget struct {
	ui.Widget
	lvItems1 *ui.Table
	lvItems2 *ui.Table

	lblInnerWidget    *ui.Label
	txtBoxInnerWidget *ui.TextBox
}

func NewTableWidget() *TableWidget {
	var c TableWidget
	c.InitWidget()
	c.SetTypeName("TableWidget")
	c.lvItems1 = ui.NewTable()
	c.lvItems1.SetColumnCount(3)
	c.lvItems1.SetColumnWidth(0, 200)
	c.lvItems1.SetColumnWidth(1, 200)
	c.lvItems1.SetColumnWidth(2, 200)
	c.lvItems1.SetColumnName(0, "Col1")
	c.lvItems1.SetColumnName(1, "Col2")
	c.lvItems1.SetColumnName(2, "Col3")
	c.lvItems1.SetHeaderRowCount(2)

	c.lvItems1.SetEditTriggerDoubleClick(true)
	c.lvItems1.SetEditTriggerEnter(true)
	c.lvItems1.SetEditTriggerF2(true)

	c.lvItems1.SetColumnCellName2(1, 0, "Col1 Header")
	c.lvItems1.SetRowCount(10)
	for i := 0; i < 10; i++ {
		c.lvItems1.SetCellText2(i, 0, "row "+fmt.Sprint(i))
		c.lvItems1.SetCellText2(i, 1, "col2 text")
		c.lvItems1.SetCellText2(i, 2, "col2 text")
	}

	c.lvItems1.SetHeaderCellSpan2(0, 0, 2, 1)
	c.lvItems1.SetHeaderCellSpan2(0, 1, 1, 2)

	c.lvItems1.SetColumnCellName2(0, 1, "SPANNED COLUMN")
	c.lvItems1.SetColumnCellName2(1, 1, "COL1 HEADER")
	c.lvItems1.SetColumnCellName2(1, 2, "COL2 HEADER")

	c.AddWidgetOnGrid(c.lvItems1, 0, 0)
	c.SetYExpandable(false)

	c.lblInnerWidget = ui.NewLabel("This is an inner widget")
	c.lblInnerWidget.SetBackgroundColor(color.RGBA{R: 90, G: 90, B: 90, A: 255})
	c.lblInnerWidget.SetAutoFillBackground(true)
	c.lvItems1.AddWidgetOnTable(c.lblInnerWidget, 1, 1, 2, 2)

	c.txtBoxInnerWidget = ui.NewTextBox()
	c.txtBoxInnerWidget.SetText("This is an inner TextBox")
	c.lvItems1.AddWidgetOnTable(c.txtBoxInnerWidget, 2, 5, 1, 1)

	c.lvItems1.SetCellImage(5, 0, TestImage(), 32)

	c.lvItems1.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		return false
	})

	c.lvItems2 = ui.NewTable()
	c.lvItems2.SetColumnCount(3)
	c.lvItems2.SetColumnWidth(0, 200)
	c.lvItems2.SetColumnWidth(1, 200)
	c.lvItems2.SetColumnWidth(2, 200)
	c.lvItems2.SetColumnName(0, "Col1")
	c.lvItems2.SetColumnName(1, "Col2")
	c.lvItems2.SetColumnName(2, "Col3")
	c.lvItems2.SetRowCount(10)
	for i := 0; i < 10; i++ {
		c.lvItems2.SetCellText2(i, 0, "row "+fmt.Sprint(i))
		c.lvItems2.SetCellText2(i, 1, "col2 text")
		c.lvItems2.SetCellText2(i, 2, "col2 text")
	}
	c.AddWidgetOnGrid(c.lvItems2, 1, 0)

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
