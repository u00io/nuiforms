package ui

func (c *Widget) AddPanel(row int, col int) *Panel {
	panel := NewPanel()
	c.AddWidget(panel, row, col)
	return panel
}

func (c *Widget) AddLabel(row int, col int, text string) *Label {
	lbl := NewLabel(text)
	c.AddWidget(lbl, row, col)
	return lbl
}

func (c *Widget) AddButton(row int, col int, text string, onPush func()) *Button {
	btn := NewButton(text)
	btn.SetOnClick(onPush)
	c.AddWidget(btn, row, col)
	return btn
}

func (c *Widget) AddHSpacer(row int, col int) *HSpacer {
	hspacer := NewHSpacer()
	c.AddWidget(hspacer, row, col)
	return hspacer
}

func (c *Widget) AddVSpacer(row int, col int) *VSpacer {
	vspacer := NewVSpacer()
	c.AddWidget(vspacer, row, col)
	return vspacer
}
