package ex00gallery

import (
	"github.com/u00io/nuiforms/ui"
)

type LayoutExample struct {
	ui.Widget
}

func NewLayoutExample() *LayoutExample {
	var c LayoutExample
	c.InitWidget()

	c.SetLayout(`
		<column>
			<row>
				<button text="Light Theme" onclick="BtnLightTheme"/>
				<button text="Dark Theme" onclick="BtnDarkTheme"/>
				<hspacer/>
			</row>
			<tabwidget>
				<tab text="Button">
					<column>
						<label text="Toolbar"/>
						<row padding="0">
							<button text="Usual Button" onclick="BtnUsualButton"/>
							<button text="Disabled Button" enabled="false" onclick="BtnDisabledButton"/>
							<hspacer/>
						</row>
					</column>
					<vspacer/>
				</tab>
				<tab text="CheckBox">
				</tab>
				<tab text="RadioButton">
				</tab>
				<tab text="ComboBox">
				</tab>
				<tab text="Layout">
				</tab>
				<tab text="Label">
				</tab>
				<tab text="NumBox">
				</tab>
				<tab text="ProgressBar">
				</tab>
				<tab text="ScrollArea">
				</tab>
				<tab text="Table">
				</tab>
				<tab text="TabWidget">
				</tab>
				<tab text="TextBox">
				</tab>
			</tabwidget>
		</column>
	`, &c, nil)

	return &c
}

func (c *LayoutExample) ClickYes() {
	println("System destroyed!")
}

func (c *LayoutExample) BtnLightTheme() {
	ui.ApplyLightTheme()
}

func (c *LayoutExample) BtnDarkTheme() {
	ui.ApplyDarkTheme()
}

func (c *LayoutExample) BtnUsualButton() {
	ui.ShowMessageBox("Message", "You clicked the button!")
}

func (c *LayoutExample) BtnDisabledButton() {
	ui.ShowMessageBox("Message", "You clicked the disabled button! How did you do that?")
}

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Layout")
	form.Panel().AddWidgetOnGrid(NewLayoutExample(), 0, 0)
}
