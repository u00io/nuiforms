package ui

type Label struct {
	Widget
	text string
}

const labelMaxWidth = 200

////////////////////////////////////////////////////////////////////////////////////
// NewLabel creates a new Label widget with the specified text.

func NewLabel(text string) *Label {
	var c Label
	c.InitWidget()
	c.SetTypeName("Label")
	//c.SetMaxWidth(labelMaxWidth)
	c.SetOnPaint(func(cnv *Canvas) {
		cnv.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignLeft, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)
	})
	c.SetText(text)
	return &c
}

/////////////////////////////////////////////////////////////////////////////////
// Label methods

func (c *Label) Text() string {
	return c.text
}

func (c *Label) SetText(text string) {
	c.text = text
	c.updateInnerSize()
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

/////////////////////////////////////////////////////////////////////////////////
// Label private methods

func (c *Label) updateInnerSize() {
	textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), c.text)
	if err != nil {
		return
	}

	c.innerHeight = textHeight
	c.innerWidth = textWidth
	if c.innerWidth > labelMaxWidth {
		c.innerWidth = labelMaxWidth
	}
	c.SetMinSize(c.innerWidth, c.innerHeight)
	c.SetMaxHeight(c.innerHeight)
}
