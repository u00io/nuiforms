package ui

import (
	"fmt"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type DialogOld struct {
	Widget

	headerPanel  *dialogHeader
	contentPanel *Panel

	title      string
	CloseEvent func()

	closed bool

	acceptButton *Button
	rejectButton *Button
	OnAccept     func()
	OnReject     func() bool

	onRejectCalled bool

	TryAccept func() bool
	OnShow    func()
}

type dialogHeader struct {
	Panel
	dialog     *DialogOld
	headerText *Label
	btnClose   *Button

	mousePointerInRect   bool
	lastMouseDownX       int
	lastMouseDownY       int
	lastMouseDownDialogX int
	lastMouseDownDialogY int
	pressed              bool
}

func NewDialogHeader() *dialogHeader {
	var c dialogHeader
	c.InitWidget()
	c.SetSize(200, 30)
	c.SetAbsolutePositioning(true)
	c.SetMaxSize(1000000, 30)
	c.SetPanelPadding(0)
	c.SetCellPadding(0)

	c.headerText = NewLabel("Header")
	c.headerText.SetPosition(0, 0)
	c.headerText.SetSize(170, 30)
	c.headerText.SetAnchors(true, true, true, true)
	c.headerText.SetOnMouseDown(c.onMouseDown)
	c.headerText.SetOnMouseUp(c.onMouseUp)
	c.headerText.SetOnMouseMove(c.onMouseMove)
	c.headerText.SetOnMouseLeave(c.onMouseLeave)
	c.AddWidget(0, 0, c.headerText)

	c.btnClose = NewButton("\u00D7")
	c.btnClose.SetPosition(170, 0)
	c.btnClose.SetSize(30, 30)
	c.btnClose.SetAnchors(false, true, true, true)
	c.btnClose.SetOnClick(func() {
		c.dialog.Reject()
	})
	c.AddWidget(0, 0, c.btnClose)
	return &c
}

func NewDialogOld(title string, width, height int) *DialogOld {
	var c DialogOld
	c.InitWidget()
	c.SetTypeName("Dialog")
	c.SetName("Dialog")
	c.SetElevation(3)
	c.SetAutoFillBackground(true)

	c.SetCloseByClickOutside(false)

	c.headerPanel = NewDialogHeader()
	c.headerPanel.dialog = &c
	c.headerPanel.SetName("DialogHeaderPanel")
	c.headerPanel.SetElevation(2)
	c.headerPanel.SetAutoFillBackground(true)
	c.AddWidget(0, 0, c.headerPanel)

	c.contentPanel = NewPanel()
	c.contentPanel.SetName("DialogContentPanel")
	c.AddWidget(1, 0, c.contentPanel)

	c.SetOnKeyDown(c.onKeyDown)

	c.SetTitle(title)
	c.Resize(width, height)

	return &c
}

func (c *dialogHeader) setTitle(title string) {
	c.headerText.SetText(title)
}

func (c *dialogHeader) MouseEnter() {
	c.mousePointerInRect = true
}

func (c *dialogHeader) MouseLeave() {
	c.mousePointerInRect = false
}

func (c *dialogHeader) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = true

	posHeaderTextX, posHeaderTextY := c.RectClientAreaOnWindow()

	x += posHeaderTextX
	y += posHeaderTextY

	c.lastMouseDownX = x
	c.lastMouseDownY = y

	c.lastMouseDownDialogX, c.lastMouseDownDialogY = c.dialog.RectClientAreaOnWindow()

	return true
}

func (c *dialogHeader) onMouseLeave() {
	if c.pressed {
		c.pressed = false
	}
}

func (c *dialogHeader) onMouseMove(x int, y int, mods nuikey.KeyModifiers) bool {
	if c.pressed {
		posHeaderTextX, posHeaderTextY := c.RectClientAreaOnWindow()

		x += posHeaderTextX
		y += posHeaderTextY

		deltaX := x - c.lastMouseDownX
		deltaY := y - c.lastMouseDownY

		c.dialog.SetPosition(c.lastMouseDownDialogX+deltaX, c.lastMouseDownDialogY+deltaY)
		//c.dialog.SetPosition(c.dialog.X()+deltaX, c.dialog.Y()+deltaY)
		fmt.Println("Dialog moved to:", c.dialog.X(), c.dialog.Y(), "Delta:", deltaX, deltaY)
		c.form.Update()
	}
	return true
}

func (c *dialogHeader) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.pressed {
		c.pressed = false
	}
	return true
}

func (c *DialogOld) Close() {
	c.Reject()
}

func (c *DialogOld) ShowDialog() {
	windowWidth := c.form.Panel().Width()
	windowHeight := c.form.Panel().Height()

	x := (windowWidth - c.Width()) / 2
	y := (windowHeight - c.Height()) / 2

	c.ShowDialogAtPos(x, y)
}

func (c *DialogOld) ShowDialogAtPos(x, y int) {
	c.SetPosition(x, y)
	c.form.Panel().AppendPopupWidget(c)
	c.ContentPanel().Focus()
	//c.Window().ProcessTabDown()

	if c.OnShow != nil {
		c.OnShow()
	}

	c.form.Update()
}

func (c *DialogOld) ContentPanel() *Panel {
	return c.contentPanel
}

func (c *DialogOld) Resize(w, h int) {
	c.SetSize(w, h)
}

func (c *DialogOld) SetAcceptButton(acceptButton *Button) {
	c.acceptButton = acceptButton
	acceptButton.SetOnClick(func() {
		c.Accept()
	})
}

func (c *DialogOld) SetRejectButton(rejectButton *Button) {
	c.rejectButton = rejectButton
	rejectButton.SetOnClick(func() {
		c.Reject()
	})
}

func (c *DialogOld) SetTitle(title string) {
	c.title = title
	c.headerPanel.setTitle(title)
}

func (c *DialogOld) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *DialogOld) Accept() {
	if c.closed {
		return
	}

	if c.TryAccept != nil {
		if !c.TryAccept() {
			return
		}
	}

	onAccept := c.OnAccept
	c.form.Panel().CloseTopPopup()
	c.closed = true
	if onAccept != nil {
		onAccept()
	}
}

func (c *DialogOld) Reject() {
	if c.closed {
		return
	}

	onReject := c.OnReject
	allowReject := true
	if onReject != nil {
		if !c.onRejectCalled {
			c.onRejectCalled = true
			allowReject = onReject()
			c.onRejectCalled = false
		} else {
			allowReject = true
		}
	}
	if allowReject {
		c.form.Panel().CloseTopPopup()
		c.closed = true
	}
}

func (c *DialogOld) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter {
		if c.acceptButton != nil {
			//c.acceptButton.Press()
			if f := c.acceptButton.GetPropFunction("onclick"); f != nil {
				f()
			}
		}
		c.Accept()
		return true
	}
	if key == nuikey.KeyEsc {
		if c.rejectButton != nil {
			//c.rejectButton.Press()
			if f := c.rejectButton.GetPropFunction("onclick"); f != nil {
				f()
			}
		}
		c.Reject()
		return true
	}
	return false
}
