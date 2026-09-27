# ToolButton

Fixed-size image button for toolbars (48x48 by default).

- Disabled state: a grayed-out copy of the image is shown and clicks are ignored.
- Checked state for toggles: a lighter background and a bright bar along the bottom.
- A lighter bottom edge, so the buttons stand out from the toolbar.

## Create

```go
img := /* load image.Image, e.g. a 32x32 icon */
btn := ui.NewToolButton(img, "Add host (A)", func() {
  // ...
})
toolbar.AddWidget(0, 0, btn)
```

## Size

```go
// The main action: twice as wide
btn.SetButtonSize(ui.ToolButtonDefaultSize*2, ui.ToolButtonDefaultSize)
```

## State

```go
btn.SetEnabled(false)       // grayed out, clicks ignored
btn.SetChecked(visible)     // the button does not toggle itself on click
```
