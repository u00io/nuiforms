package ui

type HSpacer struct {
	Widget
}

func NewHSpacer() *HSpacer {
	var c HSpacer
	c.InitWidget()
	c.SetXExpandable(true)
	return &c
}
