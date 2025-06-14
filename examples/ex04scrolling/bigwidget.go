package ex04scrolling

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func NewBigWidget() *ui.Widget {
	widget := ui.NewWidget()
	widget.SetInnerSize(1000, 1000)
	widget.SetOnPaint(func(cnv *ui.Canvas) {
		// draw a grid
		for x := 0; x < 1000; x += 50 {
			cnv.DrawLine(x, 0, x, 1000, 1, color.RGBA{200, 200, 200, 255})
		}

		for y := 0; y < 1000; y += 50 {
			cnv.DrawLine(0, y, 1000, y, 1, color.RGBA{200, 200, 200, 255})
		}

		// draw coordinates
		for x := 0; x < 1000; x += 50 {
			for y := 0; y < 1000; y += 50 {
				text := fmt.Sprint(x, "\r\n", y)
				//cnv.DrawText(x+5, y+5, text, "rotobomono", 12, color.RGBA{0, 255, 0, 255}, false)
				cnv.DrawTextMultiline(x, y, 50, 50, ui.HAlignCenter, ui.VAlignCenter, text, color.RGBA{0, 255, 0, 255}, "robotomono", 12, false)
			}
		}
	})
	widget.SetAllowScroll(true, true)
	return widget
}
