package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type ComboBox struct {
	Widget
	items         []*ComboBoxItem
	selectedIndex int
}

type ComboBoxItem struct {
	text string
	data interface{}
}

func NewComboBox() *ComboBox {
	var c ComboBox
	c.InitWidget()
	c.SetMinHeight(32)

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.OpenPopup()
			return true
		}
		return false
	})

	c.SetTypeName("ComboBox")
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetOnPaint(c.draw)
	return &c
}

func (c *ComboBox) AddItem(text string, data interface{}) {
	var item ComboBoxItem
	item.text = text
	item.data = data
	c.items = append(c.items, &item)
}

func (c *ComboBox) SetSelectedIndex(index int) {
	if index < 0 || index >= len(c.items) {
		return
	}
	c.selectedIndex = index
	UpdateMainForm()
}

func (c *ComboBox) SelectedItemText() string {
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.items) {
		return ""
	}
	return c.items[c.selectedIndex].text
}

func (c *ComboBox) SelectedItemData() interface{} {
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.items) {
		return nil
	}
	return c.items[c.selectedIndex].data
}

func (c *ComboBox) OpenPopup() {
	popup := NewComboBoxPopup()
	for _, item := range c.items {
		popup.AddItem(item.text, func(index int) {
			c.SetSelectedIndex(index)
			UpdateMainForm()
		})
	}
	x, y := c.RectClientAreaOnWindow()
	popup.ShowPopup(x, y+c.Height())
	c.SetSelectedIndex(0)
}

func (c *ComboBox) draw(cnv *Canvas) {
	backColor := c.BackgroundColorWithAddElevation(-1)
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	foreColor := c.ForegroundColor()

	itemText := c.SelectedItemText()

	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)
	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(foreColor)
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(0, 0, c.Width(), c.Height(), itemText)
}

type comboBoxPopup struct {
	Widget

	items []*comboBoxPopupItem
}

func NewComboBoxPopup() *comboBoxPopup {
	var c comboBoxPopup
	c.InitWidget()
	c.SetTypeName("ComboBoxPopup")
	c.SetAbsolutePositioning(true)
	return &c
}

func (c *comboBoxPopup) ShowPopup(x int, y int) {
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	MainForm.Panel().AppendPopupWidget(c)
	UpdateMainForm()
}

func (c *comboBoxPopup) AddItem(text string, onClick func(index int)) {
	item := newComboBoxPopupItem(len(c.items), text)
	item.parentWidgetId = c.Id()
	item.OnClick = onClick
	c.items = append(c.items, item)
	c.AddWidget(item)
}

func (c *comboBoxPopup) rebuildVisualElements() {
	yOffset := 0
	for _, item := range c.items {
		item.SetPosition(0, yOffset)
		item.SetSize(c.Width(), ContextMenuItemHeight)
		yOffset += ContextMenuItemHeight
	}
	c.SetSize(300, yOffset)
}

type comboBoxPopupItem struct {
	Widget
	index   int
	text    string
	OnClick func(index int)
}

func newComboBoxPopupItem(index int, text string) *comboBoxPopupItem {
	var item comboBoxPopupItem
	item.InitWidget()
	item.SetTypeName("ComboBoxPopupItem")
	item.SetAbsolutePositioning(true)
	item.SetMouseCursor(nuimouse.MouseCursorPointer)
	item.SetOnPaint(item.Draw)
	item.SetOnMouseDown(item.mouseDownHandler)

	item.index = index
	item.text = text
	return &item
}

func (c *comboBoxPopupItem) Draw(ctx *Canvas) {
	backColor := c.BackgroundColorWithAddElevation(-1)
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight(), backColor)
	ctx.SetHAlign(HAlignLeft)
	ctx.SetVAlign(VAlignCenter)
	ctx.SetColor(c.ForegroundColor())
	ctx.SetFontFamily(c.FontFamily())
	ctx.SetFontSize(c.FontSize())
	ctx.DrawText(0, 0, c.Width(), c.Height(), c.text)
}

func (c *comboBoxPopupItem) mouseDownHandler(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.OnClick != nil {
		c.OnClick(c.index)
	}
	MainForm.Panel().CloseTopPopup()
	return true
}
