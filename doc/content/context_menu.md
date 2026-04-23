# ContextMenu

Popup menu widget.

## Create and show

```go
menu := ui.NewContextMenu(form.Panel())
menu.AddItem("Open", func() { /* ... */ })
menu.AddItem("Quit", func() { form.Close() })
menu.ShowMenu(100, 100)
```

## Submenus

```go
sub := ui.NewContextMenu(form.Panel())
sub.AddItem("Item", func() {})
menu.AddItemWithSubmenu("More", sub)
```

