package ex12chart

import "github.com/u00io/nuiforms/ui"

type ChartExample struct {
	ui.Widget

	chart *ui.Chart
}

func NewChartExample() *ChartExample {
	var c ChartExample
	c.InitWidget()
	c.SetTypeName("ChartExample")

	c.chart = ui.NewChart()
	c.chart.SetData([]ui.ChartPoint{
		{X: 0, Y: 0},
		{X: 1, Y: 2},
		{X: 2, Y: 1},
		{X: 3, Y: 3},
	})
	c.AddWidgetOnGrid(c.chart, 0, 0)

	return &c
}

func Run(form *ui.Form) {
	chart := NewChartExample()
	form.Panel().AddWidget(chart)
}
