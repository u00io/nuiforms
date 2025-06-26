package ex03masterdetail

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type MasterWidget struct {
	ui.Widget

	panelLeft *ui.Panel
	table     *ui.Table

	panelRight *ui.Panel
	txtCol1    *ui.TextBox
	txtCol2    *ui.TextBox
	txtCol3    *ui.TextBox

	loadingDetails bool
}

func NewMasterWidget() *MasterWidget {
	var c MasterWidget
	c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 60, A: 255})
	c.InitWidget()

	c.panelLeft = ui.NewPanel()
	c.AddWidgetOnGrid(c.panelLeft, 0, 0)
	c.panelLeft.SetMinWidth(310)

	c.table = ui.NewTable()
	c.table.SetColumnCount(3)
	c.table.SetColumnName(0, "ID")
	c.table.SetColumnName(1, "Name")
	c.table.SetColumnName(2, "Description")
	c.table.SetColumnWidth(0, 50)
	c.table.SetColumnWidth(1, 100)
	c.table.SetColumnWidth(2, 100)
	c.table.SetRowCount(10)
	c.table.SetOnSelectionChanged(c.onTableSelectionChanged)
	c.panelLeft.AddWidgetOnGrid(c.table, 0, 0)

	c.panelRight = ui.NewPanel()
	c.panelRight.SetName("panelRight")
	c.AddWidgetOnGrid(c.panelRight, 1, 0)

	c.txtCol1 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol1, 0, 0)
	c.txtCol1.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(0, c.table.CurrentRow(), txt.Text())
	})

	c.txtCol2 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol2, 0, 1)
	c.txtCol2.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(1, c.table.CurrentRow(), txt.Text())
	})

	c.txtCol3 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol3, 0, 2)
	c.txtCol3.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(2, c.table.CurrentRow(), txt.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewVSpacer(), 0, 3)

	return &c
}

func (c *MasterWidget) onTableSelectionChanged(x, y int) {
	c.loadingDetails = true
	c.txtCol1.SetText(c.table.GetCellText(0, y))
	c.txtCol2.SetText(c.table.GetCellText(1, y))
	c.txtCol3.SetText(c.table.GetCellText(2, y))
	c.loadingDetails = false
}
