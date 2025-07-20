package ex15dialog

import (
	"fmt"
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type Dialog struct {
	ui.Widget

	headerPanel  *dialogHeader
	contentPanel *ui.Panel

	title      string
	CloseEvent func()

	closed bool

	acceptButton *ui.Button
	rejectButton *ui.Button
	OnAccept     func()
	OnReject     func()

	TryAccept func() bool
	OnShow    func()
}

type dialogHeader struct {
	ui.Panel
	dialog     *Dialog
	headerText *ui.Label
	btnClose   *ui.Button

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
	c.SetBackgroundColor(color.RGBA{R: 100, G: 100, B: 100, A: 255})

	c.headerText = ui.NewLabel("Header")
	c.headerText.SetPosition(0, 0)
	c.headerText.SetSize(170, 30)
	c.headerText.SetAnchors(true, true, true, true)
	c.headerText.SetOnMouseDown(c.onMouseDown)
	c.headerText.SetOnMouseUp(c.onMouseUp)
	c.headerText.SetOnMouseMove(c.onMouseMove)
	c.headerText.SetOnMouseLeave(c.onMouseLeave)
	c.AddWidget(c.headerText)

	c.btnClose = ui.NewButton("X")
	c.btnClose.SetPosition(170, 0)
	c.btnClose.SetSize(30, 30)
	c.btnClose.SetAnchors(false, true, true, true)
	c.btnClose.SetOnButtonClick(func(btn *ui.Button) {
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

	c.headerPanel = NewDialogHeader()
	c.headerPanel.dialog = &c
	c.headerPanel.SetName("DialogHeaderPanel")
	c.AddWidgetOnGrid(c.headerPanel, 0, 0)

	c.contentPanel = ui.NewPanel()
	c.contentPanel.SetName("DialogContentPanel")
	c.contentPanel.SetBackgroundColor(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	c.AddWidgetOnGrid(c.contentPanel, 0, 1)

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
		ui.UpdateMainForm()
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
	windowWidth := ui.MainForm.Panel().Width()
	windowHeight := ui.MainForm.Panel().Height()

	x := (windowWidth - c.Width()) / 2
	y := (windowHeight - c.Height()) / 2

	c.ShowDialogAtPos(x, y)
}

func (c *Dialog) ShowDialogAtPos(x, y int) {
	c.SetPosition(x, y)
	ui.MainForm.Panel().AppendPopupWidget(c)
	c.ContentPanel().Focus()
	//c.Window().ProcessTabDown()

	if c.OnShow != nil {
		c.OnShow()
	}

	ui.UpdateMainFormLayout()
}

func (c *Dialog) ContentPanel() *ui.Panel {
	return c.contentPanel
}

func (c *Dialog) Resize(w, h int) {
	c.SetSize(w, h)
}

func (c *Dialog) SetAcceptButton(acceptButton *ui.Button) {
	c.acceptButton = acceptButton
	acceptButton.SetOnButtonClick(func(btn *ui.Button) {
		c.Accept()
	})
}

func (c *Dialog) SetRejectButton(rejectButton *ui.Button) {
	c.rejectButton = rejectButton
	rejectButton.SetOnButtonClick(func(btn *ui.Button) {
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
	ui.MainForm.Panel().CloseTopPopup()
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

	ui.MainForm.Panel().CloseTopPopup()
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

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	btn := ui.NewButton("Click me")
	btn.SetOnButtonClick(func(btn *ui.Button) {
		dialog := NewDialog("Example Dialog", 400, 300)
		dialog.SetAcceptButton(ui.NewButton("Accept"))
		dialog.SetRejectButton(ui.NewButton("Reject"))
		dialog.ShowDialog()
	})
	form.Panel().AddWidgetOnGrid(btn, 0, 0)
	form.Panel().AddWidgetOnGrid(ui.NewHSpacer(), 1, 0)
	form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), 0, 10)
}
