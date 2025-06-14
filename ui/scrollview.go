package ui

import (
	"fmt"
	"image/color"
)

func NewScrollView() *Widget {
	widget := NewWidget()
	widget.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	widget.SetSize(1000, 1000)
	widget.SetAllowScroll(false, true)
	widget.innerHeight = 900 // Set a large inner height to enable scrolling

	widget.SetOnPaint(func(cnv *Canvas) {
		for i := 0; i < 33; i++ {
			cnv.DrawTextMultiline(150, i*30, widget.W()-20, 30, HAlignLeft, VAlignTop, "Line "+fmt.Sprint(i), color.RGBA{R: 0, G: 0, B: 0, A: 255}, "robotomono", 16, false)
		}
	})

	{
		// Add buttons to the scroll view
		for i := 0; i < 30; i++ {
			btn := NewButton()
			btn.SetProp("text", fmt.Sprintf("Button %d", i))
			btn.SetPosition(10, i*30)
			btn.SetSize(100, 30)
			btn.SetProp("onClick", func() {
				fmt.Println("Button clicked:", btn.GetProp("text"))
			})
			widget.AddWidget(btn)
		}

	}
	return widget
}
