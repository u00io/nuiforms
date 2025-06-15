package ui

import "image/color"

type Table struct {
	widget Widget

	rows map[int]*tableRow

	rowCount    int
	columnCount int
}

type tableRow struct {
	cells map[int]*tableCell
}

type tableCell struct {
	text string
}

func NewTable() *Table {
	var c Table
	c.widget.InitWidget()
	c.widget.SetXExpandable(true)
	c.widget.SetYExpandable(true)
	c.widget.SetOnPaint(c.draw)

	c.rows = make(map[int]*tableRow)
	return &c
}

func (c *Table) Widgeter() any {
	return &c.widget
}

func (c *Table) SetRowCount(count int) {
	c.rowCount = count
}

func (c *Table) SetColumnCount(count int) {
	c.columnCount = count
}

func (c *Table) SetCellText(row, col int, text string) {
	rowObj, exists := c.rows[row]
	if !exists {
		rowObj = &tableRow{cells: make(map[int]*tableCell)}
		c.rows[row] = rowObj
	}
	cellObj, exists := rowObj.cells[col]
	if !exists {
		cellObj = &tableCell{}
		rowObj.cells[col] = cellObj
	}
	cellObj.text = text
	UpdateMainForm()
}

func (c *Table) draw(cnv *Canvas) {
	cnv.FillRect(0, 0, c.widget.Width(), c.widget.Height(), color.RGBA{R: 255, G: 200, B: 255, A: 255})

	for rowIndex := 0; rowIndex < c.rowCount; rowIndex++ {
		rowObj, exists := c.rows[rowIndex]
		if !exists {
			continue
		}
		for colIndex := 0; colIndex < c.columnCount; colIndex++ {
			cellObj, exists := rowObj.cells[colIndex]
			if !exists {
				continue
			}
			x := colIndex * 200 // Assuming each cell is 100 pixels wide
			y := rowIndex * 30  // Assuming each cell is 30 pixels tall
			cnv.FillRect(x, y, 200, 30, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			cnv.DrawTextMultiline(x, y, 200, 30, HAlignLeft, VAlignCenter, cellObj.text, color.RGBA{R: 0, G: 0, B: 0, A: 255}, "robotomono", 16, false)
		}
	}

}
