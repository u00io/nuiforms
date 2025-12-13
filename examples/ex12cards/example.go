package ex12cards

import "github.com/u00io/nuiforms/ui"

type Cards struct {
	ui.Widget
}

func NewCards() *Cards {
	var c Cards
	c.InitWidget()
	c.SetTypeName("Cards")

	c.SetAutoFillBackground(true)
	c.SetElevation(1)

	customWidgets := make(map[string]ui.Widgeter)
	customWidgets["card1"] = NewCategoryWidget("Category 1")
	customWidgets["card2"] = NewCategoryWidget("Category 2")
	customWidgets["card3"] = NewCategoryWidget("Category 3")
	c.SetLayout(`
		<column>
			<label text="Cards Example" />
			<scrollarea>
				<widget id="card1" />
				<widget id="card2" />
				<widget id="card3" />
			</scrollarea>
			<label text="Bottom Panel" />
		</column>
	`, &c, customWidgets)

	customWidgets["card2"].SetRole("primary")

	return &c
}
