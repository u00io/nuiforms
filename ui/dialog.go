package ui

func (c *Widget) ShowDialog(title string, width int, height int, centralWidget Widgeter) {
	form := NewForm()
	form.SetTitle(title)
	form.SetAllowMinimize(false)
	form.SetAllowMaximize(false)
	form.SetSize(width, height)
	form.Panel().AddWidget(0, 0, centralWidget)
	form.ShowModal(c.form)
}
