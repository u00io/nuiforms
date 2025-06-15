package ui

type VSpacer struct {
	widget Widget
}

func NewVSpacer() *VSpacer {
	var c VSpacer
	c.widget.InitWidget()
	c.widget.SetYExpandable(true)
	return &c
}

func (c *VSpacer) Widgeter() any {
	return &c.widget
}
