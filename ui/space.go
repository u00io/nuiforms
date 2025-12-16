package ui

type Space struct {
	Widget
}

func NewSpace() *Space {
	var c Space
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetCellPadding(0)
	c.SetTypeName("Space")
	return &c
}

func (c *Space) SetSize(width, height int) {
	// Fixed size
	c.SetProp("width", width)
	c.SetProp("height", height)
	c.SetMinSize(width, height)
	c.SetMaxSize(width, height)
	c.SetXExpandable(false)
	c.SetYExpandable(false)
}

func (c *Space) MinWidth() int {
	return c.GetPropInt("width", 0)
}

func (c *Space) MinHeight() int {
	return c.GetPropInt("height", 0)
}

func (c *Space) MaxWidth() int {
	return c.GetPropInt("width", 0)
}

func (c *Space) MaxHeight() int {
	return c.GetPropInt("height", 0)
}
