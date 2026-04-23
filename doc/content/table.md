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

See `examples/ex13table` for a full demo.

