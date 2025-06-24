package ui

type VSpacer struct {
	Widget
}

func NewVSpacer() *VSpacer {
	var c VSpacer
	c.InitWidget()
	c.SetYExpandable(true)
	return &c
}
