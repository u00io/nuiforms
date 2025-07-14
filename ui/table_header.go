package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type tableHeader struct {
	Widget

	OnHeaderMouseDown func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool
	OnHeaderMouseUp   func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool
	OnHeaderMouseMove func(x int, y int, mods nuikey.KeyModifiers) nuimouse.MouseCursor
}

func newTableHeader() *tableHeader {
	var c tableHeader
	c.InitWidget()
	c.SetOnMouseDown(c.onMouseDown)
	c.SetOnMouseUp(c.onMouseUp)
	c.SetOnMouseMove(c.onMouseMove)
	//c.SetBackgroundColor(color.RGBA{R: 240, G: 240, B: 240, A: 100})
	return &c
}

func (c *tableHeader) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.OnHeaderMouseDown != nil {
		if c.OnHeaderMouseDown(button, x, y, mods) {
			return true
		}
	}
	return true
}

func (c *tableHeader) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.OnHeaderMouseUp != nil {
		if c.OnHeaderMouseUp(button, x, y, mods) {
			return true
		}
	}
	return true
}

func (c *tableHeader) onMouseMove(x int, y int, mods nuikey.KeyModifiers) bool {
	if c.OnHeaderMouseMove != nil {
		cursor := c.OnHeaderMouseMove(x, y, mods)
		c.SetMouseCursor(cursor)
	}
	return true
}
