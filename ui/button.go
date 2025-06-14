package ui

import (
	"time"

	"github.com/u00io/nui/nuimouse"
)

func NewButton() *Widget {
	c := NewWidget()
	c.SetName("Button-" + time.Now().Format("20060102150405"))
	c.SetSize(100, 30)
	c.SetOnPaint(func(cnv *Canvas) {
		backColor := GetThemeColor("background", DefaultBackground)
		if c.IsHovered() {
			backColor = GetThemeColor("button.background.hover", DefaultBackground)
			if c.GetPropBool("pressed", false) {
				backColor = GetThemeColor("button.background.pressed", DefaultBackground)
			}
		}
		cnv.FillRect(0, 0, c.W(), c.H(), backColor)
		cnv.DrawTextMultiline(0, 0, c.W(), c.H(), HAlignCenter, VAlignCenter, c.GetPropString("text", ""), GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})

	c.SetMouseCursor(nuimouse.MouseCursorPointer)

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int) {
		c.SetProp("pressed", true)
	})

	c.SetOnMouseUp(func(button nuimouse.MouseButton, x int, y int) {
		c.SetProp("pressed", false)
		if MainForm.hoverWidget == c {
			fnOnClick, ok := c.GetProp("onClick").(func())
			if ok {
				fnOnClick()
			}
		}
	})

	return c
}
