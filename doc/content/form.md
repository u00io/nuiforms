# Form

`Form` is the main application window.

## Create and run

```go
form := ui.NewForm()
form.SetTitle("My app")
form.SetSize(800, 600)
form.Exec()
```

## Key methods

- **`Panel() *ui.Panel`**: form root container.
- **`SetMainWidget(w ui.Widgeter)`**: replace the root content with a single widget.
- **`Exec()` / `ExecMaximized()`**: run the window event loop.
- **`Close()`**: close the window.

## Global events

- **`SetOnGlobalKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool)`**: intercept key presses before focused widget.

