package ui

import (
	"fmt"
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Dialog struct {
	Widget

	headerPanel  *dialogHeader
	contentPanel *Panel

	title      string
	CloseEvent func()

	closed bool

	acceptButton *Button
	rejectButton *Button
	OnAccept     func()
	OnReject     func()

	TryAccept func() bool
	OnShow    func()
}

type dialogHeader struct {
	Panel
	dialog     *Dialog
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
	c.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})

	c.headerText = NewLabel("Header")
	c.headerText.SetPosition(0, 0)
	c.headerText.SetSize(170, 30)
	c.headerText.SetAnchors(true, true, true, true)
	c.headerText.SetOnMouseDown(c.onMouseDown)
	c.headerText.SetOnMouseUp(c.onMouseUp)
	c.headerText.SetOnMouseMove(c.onMouseMove)
	c.headerText.SetOnMouseLeave(c.onMouseLeave)
	c.AddWidget(c.headerText)

	c.btnClose = NewButton("\u00D7")
	c.btnClose.SetPosition(170, 0)
	c.btnClose.SetSize(30, 30)
	c.btnClose.SetAnchors(false, true, true, true)
	c.btnClose.SetOnButtonClick(func() {
		c.dialog.Reject()
	})
	c.AddWidget(c.btnClose)
	return &c
}

func NewDialog(title string, width, height int) *Dialog {
	var c Dialog
	c.InitWidget()
	c.SetTypeName("Dialog")
	c.SetName("Dialog")
	c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	c.SetAutoFillBackground(true)

	c.headerPanel = NewDialogHeader()
	c.headerPanel.dialog = &c
	c.headerPanel.SetName("DialogHeaderPanel")
	c.headerPanel.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	c.headerPanel.SetAutoFillBackground(true)
	c.AddWidgetOnGrid(c.headerPanel, 0, 0)

	c.contentPanel = NewPanel()
	c.contentPanel.SetName("DialogContentPanel")
	c.AddWidgetOnGrid(c.contentPanel, 1, 0)

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
		UpdateMainForm()
	}
	return true
}

func (c *dialogHeader) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.pressed {
		c.pressed = false
	}
	return true
}

func (c *Dialog) Close() {
	c.Reject()
}

func (c *Dialog) ShowDialog() {
	windowWidth := MainForm.Panel().Width()
	windowHeight := MainForm.Panel().Height()

	x := (windowWidth - c.Width()) / 2
	y := (windowHeight - c.Height()) / 2

	c.ShowDialogAtPos(x, y)
}

func (c *Dialog) ShowDialogAtPos(x, y int) {
	c.SetPosition(x, y)
	MainForm.Panel().AppendPopupWidget(c)
	c.ContentPanel().Focus()
	//c.Window().ProcessTabDown()

	if c.OnShow != nil {
		c.OnShow()
	}

	UpdateMainFormLayout()
}

func (c *Dialog) ContentPanel() *Panel {
	return c.contentPanel
}

func (c *Dialog) Resize(w, h int) {
	c.SetSize(w, h)
}

func (c *Dialog) SetAcceptButton(acceptButton *Button) {
	c.acceptButton = acceptButton
	acceptButton.SetOnButtonClick(func() {
		c.Accept()
	})
}

func (c *Dialog) SetRejectButton(rejectButton *Button) {
	c.rejectButton = rejectButton
	rejectButton.SetOnButtonClick(func() {
		c.Reject()
	})
}

func (c *Dialog) SetTitle(title string) {
	c.title = title
	c.headerPanel.setTitle(title)
}

func (c *Dialog) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *Dialog) Accept() {
	if c.closed {
		return
	}

	if c.TryAccept != nil {
		if !c.TryAccept() {
			return
		}
	}

	onAccept := c.OnAccept
	MainForm.Panel().CloseTopPopup()
	c.closed = true
	if onAccept != nil {
		onAccept()
	}
}

func (c *Dialog) Reject() {
	if c.closed {
		return
	}

	onReject := c.OnReject

	MainForm.Panel().CloseTopPopup()
	c.closed = true
	if onReject != nil {
		onReject()
	}
}

func (c *Dialog) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter {
		if c.acceptButton != nil {
			c.acceptButton.Press()
		}
		c.Accept()
		return true
	}
	if key == nuikey.KeyEsc {
		if c.rejectButton != nil {
			c.rejectButton.Press()
		}
		c.Reject()
		return true
	}
	return false
}
