package ui

type HSpacer struct {
	widget Widget
}

func NewHSpacer() *HSpacer {
	var c HSpacer
	c.widget.InitWidget()
	c.widget.SetXExpandable(true)
	return &c
}

func (c *HSpacer) Widgeter() any {
	return &c.widget
}
