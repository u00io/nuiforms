package ex03masterdetail

import (
	"fmt"
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

	panelTableButtons := ui.NewPanel()
	panelTableButtons.SetName("panelTableButtons")
	panelTableButtons.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 70, A: 255})

	btnAdd := ui.NewButton()
	btnAdd.SetText("Add")
	btnAdd.SetOnButtonClick(func(btn *ui.Button) {
		row := c.table.RowCount()
		c.table.SetRowCount(row + 1)
		c.table.SetCellText(0, row, "ID"+fmt.Sprint(row))
		c.table.SetCellText(1, row, "Name"+fmt.Sprint(row))
		c.table.SetCellText(2, row, "Description"+fmt.Sprint(row))
		c.table.SetCurrentCell(0, row)
	})
	panelTableButtons.AddWidgetOnGrid(btnAdd, 0, 0)

	btnDelete := ui.NewButton()
	btnDelete.SetText("Delete")
	btnDelete.SetOnButtonClick(func(btn *ui.Button) {
	})
	panelTableButtons.AddWidgetOnGrid(btnDelete, 1, 0)

	//panelTableButtons.SetMinSize(100, 50)
	panelTableButtons.SetAllowScroll(false, false)

	c.panelLeft.AddWidgetOnGrid(panelTableButtons, 0, 0)

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
	c.panelLeft.AddWidgetOnGrid(c.table, 0, 1)

	c.panelRight = ui.NewPanel()
	c.panelRight.SetName("panelRight")
	c.AddWidgetOnGrid(c.panelRight, 1, 0)

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("ID:"), 0, 0)
	c.txtCol1 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol1, 1, 0)
	c.txtCol1.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(0, c.table.CurrentRow(), txt.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Name:"), 0, 1)
	c.txtCol2 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol2, 1, 1)
	c.txtCol2.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(1, c.table.CurrentRow(), txt.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Description:"), 0, 2)
	c.txtCol3 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol3, 1, 2)
	c.txtCol3.SetOnTextChanged(func(txt *ui.TextBox) {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText(2, c.table.CurrentRow(), txt.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Main Field:"), 0, 3)

	c.panelRight.AddWidgetOnGrid(ui.NewVSpacer(), 0, 4)

	return &c
}

func (c *MasterWidget) onTableSelectionChanged(x, y int) {
	c.loadingDetails = true
	c.txtCol1.SetText(c.table.GetCellText(0, y))
	c.txtCol2.SetText(c.table.GetCellText(1, y))
	c.txtCol3.SetText(c.table.GetCellText(2, y))
	c.loadingDetails = false
}
