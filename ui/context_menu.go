package ui

import (
	"image"
	"image/color"
)

const ContextMenuItemHeight = 32

type ContextMenu struct {
	Widget
	menuWidth  int
	menuHeight int
	items      []*ContextMenuItem
	CloseEvent func()
	closingAll func()
	parentMenu *ContextMenu
}

func NewContextMenu(parent Widgeter) *ContextMenu {
	var c ContextMenu
	c.InitWidget()
	c.SetAbsolutePositioning(true)
	c.SetName("PopupMenuPanel")
	c.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	return &c
}

func (c *ContextMenu) ShowMenu(x int, y int) {
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	MainForm.Panel().AppendPopupWidget(c)
	UpdateMainForm()
}

func (c *ContextMenu) showMenu(x int, y int, parentMenu *ContextMenu) {
	c.CloseAfterPopupWidget(parentMenu)
	c.parentMenu = parentMenu
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	//c.Window().AppendPopup(c)
	MainForm.Panel().AppendPopupWidget(c)
}

func (c *ContextMenu) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *ContextMenu) AddItem(text string, onClick func()) *ContextMenuItem {
	item := NewContextMenuItem()
	item.SetText(text)
	item.OnClick = onClick

	c.items = append(c.items, item)
	c.AddWidget(item)
	return item
}

func (c *ContextMenu) AddItemWithSubmenu(text string, img image.Image, innerMenu *ContextMenu) *ContextMenuItem {
	item := NewContextMenuItem()
	item.SetText(text)
	item.innerMenu = innerMenu
	c.items = append(c.items, item)
	c.AddWidget(item)
	return item
}

func (c *ContextMenu) RemoveAllItems() {
	c.RemoveAllWidgets()
	c.rebuildVisualElements()
	UpdateMainForm()
}

func (c *ContextMenu) OnInit() {
	c.rebuildVisualElements()
}

func (c *ContextMenu) needToClose() {
	MainForm.Panel().CloseTopPopup()
	if c.parentMenu != nil {
		c.parentMenu.needToClose()
	}
}

func (c *ContextMenu) rebuildVisualElements() {
	//c.Panel.RemoveAllWidgets()
	yOffset := 0
	for _, item := range c.items {

		//var item PopupMenuItem
		//item.InitControl(&c.Panel, &item, 0, 0, 200, PopupMenuItemHeight, 0)
		/*item.Text = itemOrig.Text
		item.OnClick = itemOrig.OnClick
		item.Image = itemOrig.Image
		item.KeyCombination = itemOrig.KeyCombination
		item.innerMenu = itemOrig.innerMenu*/
		item.needToClosePopupMenu = c.needToClose
		item.parentMenu = c

		//item.SetX(0)
		//item.SetY(yOffset)
		item.SetPosition(0, yOffset)
		//item.SetWidth(c.Width())
		//item.SetHeight(PopupMenuItemHeight)
		item.SetSize(c.Width(), ContextMenuItemHeight)
		yOffset += ContextMenuItemHeight
	}
	//c.SetWidth(300)
	//c.SetHeight(yOffset)
	c.SetSize(300, yOffset)

	c.menuWidth = 300
	c.menuHeight = yOffset
}
