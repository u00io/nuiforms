package ui

type ProgressBar struct {
	Widget
	text     string
	value    float64
	minValue float64
	maxValue float64
}

func NewProgressBar(minValue, maxValue, initValue float64) *ProgressBar {
	var c ProgressBar
	c.InitWidget()
	c.SetTypeName("ProgressBar")
	c.SetMinSize(100, 30)
	c.SetMaxSize(10000, 30)
	c.SetOnPaint(c.draw)

	c.minValue = minValue
	c.maxValue = maxValue
	c.value = initValue

	return &c
}

func (c *ProgressBar) Text() string {
	return c.text
}

func (c *ProgressBar) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *ProgressBar) SetValue(value float64) {
	if value < c.minValue {
		value = c.minValue
	}
	if value > c.maxValue {
		value = c.maxValue
	}
	if c.value == value {
		return
	}
	c.value = value
	UpdateMainForm()
}

func (c *ProgressBar) Value() float64 {
	return c.value
}

func (c *ProgressBar) SetMinValue(minValue float64) {
	c.minValue = minValue
	UpdateMainForm()
}

func (c *ProgressBar) SetMaxValue(maxValue float64) {
	c.maxValue = maxValue
	UpdateMainForm()
}

func (c *ProgressBar) draw(cnv *Canvas) {
	percents := (c.value - c.minValue) / (c.maxValue - c.minValue)
	if percents < 0 {
		percents = 0
	}
	if percents > 1 {
		percents = 1
	}

	padding := 2
	effectiveWidth := c.Width() - padding*2
	effectiveHeight := c.Height() - padding*2

	cnv.SetColor(c.ForegroundColor())
	cnv.FillRect(padding, padding, int(float64(effectiveWidth)*percents), effectiveHeight, c.ForegroundColor())

	if len(c.text) > 0 {
		cnv.SetHAlign(HAlignCenter)
		cnv.SetVAlign(VAlignCenter)
		cnv.SetColor(c.ForegroundColor())
		cnv.SetFontFamily(c.FontFamily())
		cnv.SetFontSize(c.FontSize())
		cnv.DrawText(0, 0, c.Width(), c.Height(), c.text)
	}

	cnv.SetColor(c.ForegroundColor())
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}
