package ui

type Dialoger interface {
	onDialogShow()
	onDialogReject() bool
}

type DialogContent struct {
	Widget

	OnDialogShow   func()
	OnDialogReject func() bool
}

func (c *DialogContent) InitWidget() {
	c.Widget.InitWidget()
}

func (c *DialogContent) onDialogShow() {
	if c.OnDialogShow != nil {
		c.OnDialogShow()
	}
}

func (c *DialogContent) onDialogReject() bool {
	if c.OnDialogReject != nil {
		return c.OnDialogReject()
	}
	return true
}

func (c *Widget) ShowDialog(title string, width int, height int, centralWidget Widgeter) {
	form := NewForm()
	form.SetTitle(title)
	form.SetAllowMinimize(false)
	form.SetAllowMaximize(false)
	form.SetSize(width, height)
	form.Panel().AddWidget(0, 0, centralWidget)
	form.ShowModal(c.form)

	if w, ok := centralWidget.(Dialoger); ok {
		w.onDialogShow()
	}

	form.OnClose = func() bool {
		if w, ok := centralWidget.(Dialoger); ok {
			return w.onDialogReject()
		}
		return true
	}
}
