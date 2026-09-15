package ex00gallery

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

type ExampleTableSelection struct {
	ui.Widget

	lblStatus *ui.Label
	table     *ui.Table
}

func NewExampleTableSelection() *ExampleTableSelection {
	var c ExampleTableSelection
	c.InitWidget()

	c.lblStatus = c.AddLabel(0, 0, "Click a cell to select it")
	c.lblStatus.SetUnderline(true)

	cbRowSelect := ui.NewCheckbox("Highlight whole row")
	cbRowSelect.SetChecked(true)
	c.AddWidget(1, 0, cbRowSelect)

	cbCellSelect := ui.NewCheckbox("Highlight individual cell")
	cbCellSelect.SetChecked(true)
	c.AddWidget(1, 1, cbCellSelect)

	c.table = ui.NewTable()
	c.table.SetColumnCount(3)
	c.table.SetColumnName(0, "Product")
	c.table.SetColumnName(1, "Price")
	c.table.SetColumnName(2, "In Stock")
	c.table.SetColumnWidth(0, 200)
	c.table.SetColumnWidth(1, 120)
	c.table.SetColumnWidth(2, 120)

	cbRowSelect.SetOnStateChanged(func() { c.table.SetSelectingRow(cbRowSelect.Checked()) })
	cbCellSelect.SetOnStateChanged(func() { c.table.SetSelectingCell(cbCellSelect.Checked()) })

	rows := [][3]string{
		{"Keyboard", "49.99", "Yes"},
		{"Mouse", "19.99", "Yes"},
		{"Monitor", "199.99", "No"},
		{"Webcam", "39.99", "Yes"},
	}
	c.table.SetRowCount(len(rows))
	for row, r := range rows {
		c.table.SetCellText2(row, 0, r[0])
		c.table.SetCellText2(row, 1, r[1])
		c.table.SetCellText2(row, 2, r[2])
	}

	c.table.SetOnSelectionChanged(func(row, col int) {
		c.setStatus(fmt.Sprintf("Selected row %d, column %d (%s)", row, col, c.table.ColumnName(col)))
	})

	c.AddWidget(2, 0, c.table)

	return &c
}

func (c *ExampleTableSelection) setStatus(text string) {
	c.lblStatus.SetText(text)
}
