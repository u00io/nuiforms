package ui

func NewTextBlock() *Widget {
	c := NewWidget()
	c.SetSize(100, 20)
	c.SetOnPaint(func(cnv *Canvas) {
		cnv.SetColor(DefaultBackground)
		cnv.FillRect(0, 0, c.W(), c.H(), DefaultBackground)
		cnv.DrawTextMultiline(0, 0, c.W(), c.H(), HAlignLeft, VAlignTop, c.GetPropString("text", ""), GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	return c
}
