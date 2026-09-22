# Table

Grid/table widget with selection and editing support.

## Create

```go
t := ui.NewTable()
```

## Common callbacks

- `SetOnCellChanged(func(row, col int, text string, data interface{}) bool)`
- `SetOnSelectionChanged(func(x, y int))`
- `SetOnColumnClick(func(col int))`

## Selection

- `SetSelectingRows(bool)` - selection unit: whole rows (`true`, the default) or individual cells (`false`).
- `SetMultiselect(bool)` - allow selecting more than one row/cell. When enabled, standard multi-select gestures work: drag, Shift+click (range from the anchor), Ctrl+click (toggle one item), Ctrl+Shift+click (add a range), and Ctrl+A (select all).
- `SelectedRows() []int` / `SelectedCells() []TableCellPos` - the current selection.
- `IsRowSelected(row int) bool` / `IsCellSelected(row, col int) bool`

See `examples/ex00gallery` (Table page, "Selection" tab) for a full demo.

