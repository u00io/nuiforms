package ui

type Panel struct {
	widget Widget
}

func NewPanel() *Panel {
	var c Panel
	c.widget.InitWidget()
	c.widget.SetXExpandable(false)
	c.widget.SetYExpandable(false)
	return &c
}

func (c *Panel) Widgeter() any {
	return &c.widget
}

func (c *Panel) SetName(name string) {
	c.widget.SetName(name)
}

func (c *Panel) AddWidgetOnGrid(w any, x int, y int) {
	c.widget.AddWidgetOnGrid(w, x, y)
}

func (c *Panel) AddWidget(w any) {
	c.widget.AddWidget(w)
}

func (c *Panel) RemoveWidget(w any) {
	c.widget.RemoveWidget(w)
}

func (c *Panel) SetPosition(x, y int) {
	c.widget.SetPosition(x, y)
}

func (c *Panel) SetSize(w, h int) {
	c.widget.SetSize(w, h)
}

func (c *Panel) SetPanelPadding(padding int) {
	c.widget.SetPanelPadding(padding)
}

func (c *Panel) SetAnchors(left, right, top, bottom bool) {
	c.widget.SetAnchors(left, right, top, bottom)
}

func (c *Panel) SetAbsolutePositioning(absolute bool) {
	c.widget.SetAbsolutePositioning(absolute)
}

func (c *Panel) IsVisible() bool {
	return c.widget.IsVisible()
}

func (c *Panel) Width() int {
	return c.widget.Width()
}

func (c *Panel) Height() int {
	return c.widget.Height()
}

func (c *Panel) SetXExpandable(expandable bool) {
	c.widget.SetXExpandable(expandable)
}

func (c *Panel) SetYExpandable(expandable bool) {
	c.widget.SetYExpandable(expandable)
}
