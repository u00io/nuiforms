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
	//c.panelLeft.SetMinWidth(310)

	panelTableButtons := ui.NewPanel()
	panelTableButtons.SetName("panelTableButtons")
	panelTableButtons.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 70, A: 255})

	btnAdd := ui.NewButton("Add")
	btnAdd.SetOnClick(func() {
		row := c.table.RowCount()
		c.table.SetRowCount(row + 1)
		c.table.SetCellText2(row, 0, "ID"+fmt.Sprint(row))
		c.table.SetCellText2(row, 1, "Name"+fmt.Sprint(row))
		c.table.SetCellText2(row, 2, "Description"+fmt.Sprint(row))
		c.table.SetCurrentCell2(row, 0)
	})
	panelTableButtons.AddWidgetOnGrid(btnAdd, 0, 0)

	btnDelete := ui.NewButton("Delete")
	btnDelete.SetOnClick(func() {
		selectedRow := c.table.CurrentRow()
		if selectedRow < 0 || selectedRow >= c.table.RowCount() {
			return
		}
		// copy the next rows up
		for i := selectedRow; i < c.table.RowCount()-1; i++ {
			c.table.SetCellText2(i, 0, c.table.GetCellText2(i+1, 0))
			c.table.SetCellText2(i, 1, c.table.GetCellText2(i+1, 1))
			c.table.SetCellText2(i, 2, c.table.GetCellText2(i+1, 2))
		}
		// remove the last row
		c.table.SetRowCount(c.table.RowCount() - 1)
		if c.table.RowCount() > 0 {
			if selectedRow >= c.table.RowCount() {
				c.table.SetCurrentCell2(c.table.RowCount()-1, 0)
			} else {
				c.table.SetCurrentCell2(selectedRow, 0)
			}
		}
	})
	panelTableButtons.AddWidgetOnGrid(btnDelete, 0, 1)

	//panelTableButtons.SetMinSize(100, 50)
	//panelTableButtons.SetAllowScroll(false, false)

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
	c.panelLeft.AddWidgetOnGrid(c.table, 1, 0)

	c.panelRight = ui.NewPanel()
	c.panelRight.SetName("panelRight")
	c.AddWidgetOnGrid(c.panelRight, 0, 1)

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("ID:"), 0, 0)
	c.txtCol1 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol1, 0, 1)
	c.txtCol1.SetOnTextChanged(func() {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText2(c.table.CurrentRow(), 0, c.txtCol1.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Name:"), 1, 0)
	c.txtCol2 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol2, 1, 1)
	c.txtCol2.SetOnTextChanged(func() {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText2(c.table.CurrentRow(), 1, c.txtCol2.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Description:"), 2, 0)
	c.txtCol3 = ui.NewTextBox()
	c.panelRight.AddWidgetOnGrid(c.txtCol3, 2, 1)
	c.txtCol3.SetOnTextChanged(func() {
		if c.loadingDetails {
			return
		}
		c.table.SetCellText2(c.table.CurrentRow(), 2, c.txtCol3.Text())
	})

	c.panelRight.AddWidgetOnGrid(ui.NewLabel("Main Field:"), 3, 0)
	c.panelRight.AddWidgetOnGrid(ui.NewVSpacer(), 4, 0)

	return &c
}

func (c *MasterWidget) onTableSelectionChanged(x, y int) {
	c.loadingDetails = true
	c.txtCol1.SetText(c.table.GetCellText2(y, 0))
	c.txtCol2.SetText(c.table.GetCellText2(y, 1))
	c.txtCol3.SetText(c.table.GetCellText2(y, 2))
	c.loadingDetails = false
}
