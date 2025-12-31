package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Label struct {
	Widget
}

const labelMaxWidth = 1500

// //////////////////////////////////////////////////////////////////////////////////
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
		cnv.SetColor(c.ForegroundColor())
		cnv.DrawText(0, 0, c.Width(), c.Height(), c.GetPropString("text", ""))

		if c.GetPropBool("underline", false) {
			cnv.SetColor(c.ForegroundColor())
			widgetWidth := c.Width()
			textHeight := c.innerHeight
			cnv.DrawLine(0, textHeight-2, widgetWidth, textHeight-2, 1, c.ForegroundColor())
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

// ///////////////////////////////////////////////////////////////////////////////
// Props
func (c *Label) ProcessPropChange(key string, value interface{}) {
	c.updateInnerSize()
	UpdateMainFormLayout()
}

/////////////////////////////////////////////////////////////////////////////////
// Label methods

// Text returns the text of the Label.
func (c *Label) Text() string {
	return c.GetPropString("text", "")
}

// SetText sets the text of the Label.
func (c *Label) SetText(text string) {
	c.SetProp("text", text)
	c.updateInnerSize()
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

// TextAlign returns the horizontal alignment of the Label's text.
func (c *Label) TextAlign() HAlign {
	return c.GetHAlign("textAlign", HAlignLeft)
}

// SetTextAlign sets the horizontal alignment of the Label's text.
func (c *Label) SetTextAlign(align HAlign) {
	c.SetProp("textAlign", align.String())
	c.updateInnerSize()
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

// IsUnderline returns whether the Label's text is underlined.
func (c *Label) IsUnderline() bool {
	return c.GetPropBool("underline", false)
}

// SetUnderline sets whether the Label's text is underlined.
func (c *Label) SetUnderline(underline bool) {
	c.SetProp("underline", underline)
	if MainForm != nil {
		MainForm.Update()
	}
	UpdateMainFormLayout()
}

/////////////////////////////////////////////////////////////////////////////////
// Label private methods

// updateInnerSize updates the inner size of the Label based on its text.
func (c *Label) updateInnerSize() {
	textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), c.Text())
	if err != nil {
		return
	}

	// Ensure a minimum height for the label
	if textHeight < DefaultUiLineHeight {
		textHeight = DefaultUiLineHeight
	}

	c.innerHeight = textHeight
	c.innerWidth = textWidth
	if c.innerWidth > labelMaxWidth {
		c.innerWidth = labelMaxWidth
	}
	c.SetMinSize(c.innerWidth, c.innerHeight)
	//fmt.Println("Label Min Width:", c.innerWidth, c.innerHeight)
	//c.SetMinSize(300, c.innerHeight)
	c.SetMaxHeight(c.innerHeight)
}
