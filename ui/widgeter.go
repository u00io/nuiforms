package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Widgeter interface {
	Id() string
	ParentWidgetId() string
	SetParentWidgetId(id string)
	FullPath() []string
	TypeName() string
	Name() string
	X() int
	Y() int
	Width() int
	Height() int
	InnerWidth() int
	InnerHeight() int

	Widgets() []Widgeter

	SetName(name string)
	SetPosition(x, y int)
	SetSize(width, height int)
	SetAnchors(left, top, right, bottom bool)

	getWidgetAt(x, y int) Widgeter
	findWidgetAt(x, y int) Widgeter
	Focus()

	ProcessPaint(cnv *Canvas)
	ProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers, allowTranslateToChildren bool) bool
	ProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers, onlyForWidgetId string, allowTranslateToChildren bool) bool
	ProcessMouseMove(x int, y int, mods nuikey.KeyModifiers, allowTranslateToChildren bool) bool
	ProcessMouseLeave() bool
	ProcessMouseEnter() bool
	ProcessKeyDown(keyCode nuikey.Key, mods nuikey.KeyModifiers) bool
	ProcessKeyUp(keyCode nuikey.Key, mods nuikey.KeyModifiers) bool
	ProcessMouseDblClick(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool
	ProcessChar(char rune, mods nuikey.KeyModifiers) bool
	ProcessMouseWheel(deltaX int, deltaY int) bool
	ProcessTimer()

	SetMouseCursor(cursor nuimouse.MouseCursor)
	MouseCursor() nuimouse.MouseCursor

	Anchors() (left, top, right, bottom bool)

	AddWidget(widget Widgeter)
	AddWidgetOnGrid(widget Widgeter, gridX, gridY int)
	RemoveWidget(widget Widgeter)

	IsVisible() bool

	GridX() int
	GridY() int
	SetGridPosition(x, y int)

	XExpandable() bool
	YExpandable() bool

	MinWidth() int
	MinHeight() int

	MaxWidth() int
	MaxHeight() int

	SetAbsolutePositioning(absolute bool)
	SetXExpandable(expandable bool)
	SetYExpandable(expandable bool)
}
