package ui

import "image/color"

// timeChartTheme holds every color the TimeChart draws with. The chart picks
// the dark or light one on each paint from IsDarkTheme, so switching the
// application theme only needs a repaint.
type timeChartTheme struct {
	background color.Color
	border     color.Color
	axis       color.Color
	text       color.Color
	grid       color.Color

	// Automatic series colors, in the order series are added.
	palette []color.Color

	candleUp   color.Color
	candleDown color.Color

	marker        color.Color
	markerLabelBg color.Color

	selection       color.Color
	selectionDrag   color.Color
	selectionEdge   color.Color
	selectionHeader color.Color
	headerText      color.Color
	closeHover      color.Color
	zoomBand        color.Color

	badHatchAlpha uint8
}

var timeChartDarkTheme = timeChartTheme{
	background: ColorFromHex("#0d1117"),
	border:     ColorFromHex("#30363d"),
	axis:       ColorFromHex("#6e7681"),
	text:       ColorFromHex("#8b949e"),
	grid:       ColorFromHex("#21262d"),
	palette: []color.Color{
		ColorFromHex("#4fc3f7"),
		ColorFromHex("#f0883e"),
		ColorFromHex("#a371f7"),
		ColorFromHex("#3fb950"),
		ColorFromHex("#f778ba"),
		ColorFromHex("#d29922"),
	},
	candleUp:        ColorFromHex("#3fb950"),
	candleDown:      ColorFromHex("#f85149"),
	marker:          ColorFromHex("#e3b341"),
	markerLabelBg:   timeChartAlpha("#0d1117", 200),
	selection:       timeChartAlpha("#58a6ff", 50),
	selectionDrag:   timeChartAlpha("#58a6ff", 70),
	selectionEdge:   timeChartAlpha("#58a6ff", 220),
	selectionHeader: timeChartAlpha("#1f6feb", 170),
	headerText:      ColorFromHex("#ffffff"),
	closeHover:      timeChartAlpha("#ffffff", 60),
	zoomBand:        timeChartAlpha("#ffffff", 40),
	badHatchAlpha:   80,
}

var timeChartLightTheme = timeChartTheme{
	background: ColorFromHex("#ffffff"),
	border:     ColorFromHex("#d0d7de"),
	axis:       ColorFromHex("#8c959f"),
	text:       ColorFromHex("#57606a"),
	grid:       ColorFromHex("#eaeef2"),
	palette: []color.Color{
		ColorFromHex("#0969da"),
		ColorFromHex("#d1570a"),
		ColorFromHex("#8250df"),
		ColorFromHex("#1a7f37"),
		ColorFromHex("#bf3989"),
		ColorFromHex("#9a6700"),
	},
	candleUp:        ColorFromHex("#1a7f37"),
	candleDown:      ColorFromHex("#cf222e"),
	marker:          ColorFromHex("#9a6700"),
	markerLabelBg:   timeChartAlpha("#ffffff", 220),
	selection:       timeChartAlpha("#0969da", 30),
	selectionDrag:   timeChartAlpha("#0969da", 50),
	selectionEdge:   timeChartAlpha("#0969da", 200),
	selectionHeader: timeChartAlpha("#0969da", 190),
	headerText:      ColorFromHex("#ffffff"),
	closeHover:      timeChartAlpha("#ffffff", 70),
	zoomBand:        timeChartAlpha("#000000", 25),
	badHatchAlpha:   90,
}

func currentTimeChartTheme() *timeChartTheme {
	if IsDarkTheme {
		return &timeChartDarkTheme
	}
	return &timeChartLightTheme
}

// timeChartAlpha returns a translucent color for Canvas.MixPixel (used by
// FillRect and axis-aligned DrawLine), which blends the RGB channels as
// given: they must be straight, not premultiplied as color.RGBA normally is.
// A color.NRGBA would be premultiplied by its RGBA() method and come out
// darker than intended.
func timeChartAlpha(hex string, a uint8) color.RGBA {
	c := ColorFromHex(hex)
	c.A = a
	return c
}

// timeChartStraightAlpha is timeChartAlpha for an arbitrary color.
func timeChartStraightAlpha(col color.Color, a uint8) color.RGBA {
	n := color.NRGBAModel.Convert(col).(color.NRGBA)
	return color.RGBA{R: n.R, G: n.G, B: n.B, A: a}
}
