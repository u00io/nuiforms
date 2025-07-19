package ex14popup

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	btn := ui.NewButton("Click me")
	btn.SetOnButtonClick(func(btn *ui.Button) {
		lbl := ui.NewLabel("This is a popup message!")
		lbl.SetPosition(100, 100)
		lbl.SetSize(300, 300)
		lbl.SetBackgroundColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		form.Panel().AppendPopupWidget(lbl)
	})
	form.Panel().AddWidgetOnGrid(btn, 0, 0)
	form.Panel().AddWidgetOnGrid(ui.NewHSpacer(), 1, 0)

	lbl := ui.NewLabel("Right-click for context menu")
	form.Panel().AddWidgetOnGrid(lbl, 0, 1)

	form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), 0, 2)

	contextMenu := ui.NewContextMenu(nil)
	contextMenu.AddItem("Item 1", func() {
		fmt.Println("Item 1 clicked")
	})
	contextMenu.AddItem("Item 2", func() {
		fmt.Println("Item 2 clicked")
	})

	lbl.SetContextMenu(contextMenu)
}
