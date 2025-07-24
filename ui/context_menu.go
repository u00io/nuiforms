package ui

const ContextMenuItemHeight = 32

type ContextMenu struct {
	Widget
	menuWidth  int
	menuHeight int
	items      []*ContextMenuItem
	CloseEvent func()
	parentMenu *ContextMenu
}

func NewContextMenu(parent Widgeter) *ContextMenu {
	var c ContextMenu
	c.InitWidget()
	c.SetAbsolutePositioning(true)
	c.SetName("PopupMenuPanel")
	c.SetBackgroundColor(c.BackgroundColor())
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
	item.parentWidgetId = c.Id()
	item.SetText(text)
	item.OnClick = onClick

	c.items = append(c.items, item)
	c.AddWidget(item)
	return item
}

func (c *ContextMenu) AddItemWithSubmenu(text string, innerMenu *ContextMenu) *ContextMenuItem {
	item := NewContextMenuItem()
	item.parentWidgetId = c.Id()
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
	yOffset := 0
	for _, item := range c.items {
		item.needToClosePopupMenu = c.needToClose
		item.parentMenu = c
		item.SetPosition(0, yOffset)
		item.SetSize(c.Width(), ContextMenuItemHeight)
		yOffset += ContextMenuItemHeight
	}
	c.SetSize(300, yOffset)
	c.menuWidth = 300
	c.menuHeight = yOffset
}
