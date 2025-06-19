package ex08dialog

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	panel := form.Panel()

	panel1 := ui.NewPanel()
	panel1.SetPanelPadding(10)
	panel.AddWidgetOnGrid(panel1, 0, 0)
	txtBox := ui.NewTextBox()
	panel1.AddWidgetOnGrid(txtBox, 0, 0)

	panel2 := ui.NewPanel()
	panel.AddWidgetOnGrid(panel2, 0, 1)

	//panel2.AddWidgetOnGrid(ui.NewVSpacer(), 0, 0)
	table := ui.NewTable()

	colCount := 5
	rowCount := 50

	table.SetRowCount(rowCount)
	table.SetColumnCount(colCount)

	for i := 0; i < colCount; i++ {
		table.SetColumnName(i, "Column "+fmt.Sprint(i+1))
		table.SetColumnWidth(i, 100)
	}

	for r := 0; r < rowCount; r++ {
		for column := 0; column < colCount; column++ {
			table.SetCellText(column, r, fmt.Sprintf("%dx%d", r+1, column+1))
		}
	}

	panel2.AddWidgetOnGrid(table, 0, 0)

	panel3 := ui.NewPanel()
	panel.AddWidgetOnGrid(panel3, 0, 2)
	//panel3.AddWidgetOnGrid(ui.NewHSpacer(), 0, 0)
	panel3.AddWidgetOnGrid(ui.NewTextBox(), 0, 0)
	btnOK := ui.NewButton()
	btnOK.SetText("OK")
	panel3.AddWidgetOnGrid(btnOK, 1, 0)
	btnCancel := ui.NewButton()
	btnCancel.SetText("Cancel")
	btnCancel.SetOnButtonClick(func(btn *ui.Button) {
		ui.MainForm.Close()
	})
	panel3.AddWidgetOnGrid(btnCancel, 2, 0)

	form.Exec()
}
