package ui

type Panel struct {
	Widget
}

func NewPanel() *Panel {
	var c Panel
	c.InitWidget()
	c.SetName(c.Name() + "-Panel")
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	c.SetAllowScroll(true, true)
	return &c
}
