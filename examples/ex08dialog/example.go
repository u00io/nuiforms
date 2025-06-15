package ex08dialog

import "github.com/u00io/nuiforms/ui"

func Run() {
	form := ui.NewForm()
	panel := form.Panel()

	panel1 := ui.NewPanel()
	panel.AddWidgetOnGrid(panel1, 0, 0)
	txtBox := ui.NewTextBox()
	panel1.AddWidgetOnGrid(txtBox, 0, 0)

	panel2 := ui.NewPanel()
	panel.AddWidgetOnGrid(panel2, 0, 1)
	table := ui.NewTable()

	table.SetRowCount(3)
	table.SetColumnCount(2)

	table.SetCellText(0, 0, "Row 1, Col 1")
	table.SetCellText(0, 1, "Row 1, Col 2")
	table.SetCellText(1, 0, "Row 2, Col 1")
	table.SetCellText(1, 1, "Row 2, Col 2")
	table.SetCellText(2, 0, "Row 3, Col 1")
	table.SetCellText(2, 1, "Row 3, Col 2")

	panel2.AddWidgetOnGrid(table, 0, 0)

	panel3 := ui.NewPanel()
	panel.AddWidgetOnGrid(panel3, 0, 2)
	panel3.AddWidgetOnGrid(ui.NewHSpacer(), 0, 0)
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
