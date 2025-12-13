package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Label struct {
	Widget
	//text string
	//underline bool
	//textAlign HAlign
}

const labelMaxWidth = 1500

////////////////////////////////////////////////////////////////////////////////////
// NewLabel creates a new Label widget with the specified text.

func NewLabel(text string) *Label {
	var c Label
	c.InitWidget()
	c.SetTypeName("Label")
	c.SetOnPaint(func(cnv *Canvas) {
		cnv.SetHAlign(c.GetHAlign("textAlign", HAlignLeft))
		cnv.SetVAlign(VAlignCenter)
		cnv.SetFontFamily(c.FontFamily())
		cnv.SetFontSize(c.FontSize())
		cnv.SetColor(c.Color())
		cnv.DrawText(0, 0, c.Width(), c.Height(), c.GetPropString("text", ""))

		if c.GetPropBool("underline", false) {
			cnv.SetColor(c.Color())
			widgetWidth := c.Width()
			textHeight := c.innerHeight
			cnv.DrawLine(0, textHeight-2, widgetWidth, textHeight-2, 1, c.Color())
		}
	})
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		// Labels do not handle mouse down events by default
		return true
	})

	c.SetOnMouseMove(func(x int, y int, mods nuikey.KeyModifiers) bool {
		return true
	})
	c.SetText(text)
	return &c
}

/////////////////////////////////////////////////////////////////////////////////
// Label methods

func (c *Label) Text() string {
	return c.GetPropString("text", "")
}

func (c *Label) SetText(text string) {
	c.SetProp("text", text)
	c.updateInnerSize()
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

func (c *Label) TextAlign() HAlign {
	return c.GetHAlign("textAlign", HAlignLeft)
}

func (c *Label) SetTextAlign(align HAlign) {
	c.SetProp("textAlign", align.String())
	c.updateInnerSize()
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

func (c *Label) IsUnderline() bool {
	return c.GetPropBool("underline", false)
}

func (c *Label) SetUnderline(underline bool) {
	c.SetProp("underline", underline)
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

/////////////////////////////////////////////////////////////////////////////////
// Label private methods

func (c *Label) updateInnerSize() {
	textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), c.Text())
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
