package ex06font

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	panel := ui.NewPanel()

	progressBar := ui.NewProgressBar(0, 100, 0)
	cb1 := ui.NewCheckbox("Checkbox 1")
	cb1.SetOnStateChanged(func(btn *ui.Checkbox, checked bool) {
		progressBar.SetValue(progressBar.Value() + 1)
		progressBar.SetText(fmt.Sprintf("Progress: %.0f%%", progressBar.Value()))
	})

	rb1 := ui.NewRadioButton("Option 1")
	rb1.SetOnStateChanged(func(btn *ui.RadioButton, checked bool) {
		fmt.Println("RadioButton 1 checked:", checked)
		cb1.SetChecked(checked)
	})
	panel.AddWidgetOnGrid(rb1, 0, 0)

	rb2 := ui.NewRadioButton("Option 2")
	rb2.SetOnStateChanged(func(btn *ui.RadioButton, checked bool) {
		fmt.Println("RadioButton 2 checked:", checked)
	})
	panel.AddWidgetOnGrid(rb2, 0, 1)

	rb1.SetChecked(true) // Set the first radio button as checked by default

	panel.AddWidgetOnGrid(cb1, 0, 2)

	panel.AddWidgetOnGrid(progressBar, 0, 3)

	form.Panel().AddWidgetOnGrid(panel, 0, 0)
}
