package ui

import (
	"github.com/u00io/nui/nuimouse"
)

type ChartPoint struct {
	X float64
	Y float64
}

type Chart struct {
	Widget
	points []ChartPoint
}

func NewChart() *Chart {
	var c Chart
	c.InitWidget()
	c.SetTypeName("Chart")
	c.SetMinSize(DefaultButtonMinWidth, DefaultUiLineHeight)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetCanBeFocused(true)
	c.SetElevation(1)

	c.SetOnPaint(c.draw)
	c.SetOnKeyDown(c.onKeyDown)

	c.SetProp("padding", 6)

	c.SetXExpandable(true)
	c.SetYExpandable(true)

	return &c
}

func (c *Chart) Data() []ChartPoint {
	return c.points
}

func (c *Chart) SetData(points []ChartPoint) {
	c.points = points
	UpdateMainForm()
}

func (c *Chart) draw(cnv *Canvas) {
	// fill background
	cnv.FillRect(0, 0, c.Width(), c.Height(), ColorFromHex("#000000"))

	// fill border
	cnv.SetColor(ColorFromHex("#00AA55"))
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}
