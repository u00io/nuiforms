package ex12chart

import (
	"math"
	"math/rand"

	"github.com/u00io/nuiforms/ui"
)

const (
	samplePoints = 200
)

type ChartExample struct {
	ui.Widget

	chart *ui.Chart
}

// applyDataset sets the Y-axis title and chart points for a preset.
func (c *ChartExample) applyDataset(ylabel string, gen func() []ui.ChartPoint) {
	c.chart.SetProp("ylabel", ylabel)
	c.chart.SetData(gen())
}

func (c *ChartExample) wireButton(btn *ui.Button, ylabel string, gen func() []ui.ChartPoint) {
	btn.SetPropFunction("onclick", func() { c.applyDataset(ylabel, gen) })
}

func newChartWithButtons(c *ChartExample) {
	c.chart = ui.NewChart()
	c.chart.SetProp("xlabel", "t")
	c.applyDataset("sin(t)", dataSine)
	// Toolbar: one cell in the example grid (one column) so the chart below gets full width.
	bar := ui.NewPanel()
	presets := []struct {
		title  string
		ylabel string
		gen    func() []ui.ChartPoint
	}{
		{"Sine", "sin(t)", dataSine},
		{"Two sines", "sin(t) + 0.6*sin(3t)", dataTwoSines},
		{"Noise", "gaussian noise (seeded)", dataNoise},
		{"Spikes", "sine + narrow pulses", dataSpikes},
		{"Damped", "exp(-0.1t)*sin(3t)", dataDamped},
		{"Step-ish", "floor(t) mod 5 + 0.1*sin(20t)", dataSteps},
	}
	for i, p := range presets {
		b := ui.NewButton(p.title)
		c.wireButton(b, p.ylabel, p.gen)
		bar.AddWidgetOnGrid(b, 0, i) // one row, columns 0..n-1
	}
	c.AddWidgetOnGrid(bar, 0, 0)      // first row: button strip
	c.AddWidgetOnGrid(c.chart, 1, 0)  // second row: chart
}

func NewChartExample() *ChartExample {
	var c ChartExample
	c.InitWidget()
	c.SetTypeName("ChartExample")
	newChartWithButtons(&c)
	return &c
}

func Run(form *ui.Form) {
	chart := NewChartExample()
	form.Panel().AddWidget(chart)
}

// dataSine: single period over [0, 2π).
func dataSine() []ui.ChartPoint {
	return sampleSeries(0, 2*math.Pi, samplePoints, func(t float64) float64 {
		return math.Sin(t)
	})
}

// dataTwoSines: two harmonics (as in the original example).
func dataTwoSines() []ui.ChartPoint {
	return sampleSeries(0, 4*math.Pi, samplePoints, func(t float64) float64 {
		return math.Sin(t) + 0.6*math.Sin(3*t)
	})
}

// dataNoise: pseudo-random walk style noise, deterministic seed for stable redraws.
func dataNoise() []ui.ChartPoint {
	rng := rand.New(rand.NewSource(42))
	pts := make([]ui.ChartPoint, samplePoints)
	t0, t1 := 0.0, 8.0
	var y float64
	for i := 0; i < samplePoints; i++ {
		t := t0 + float64(i)/float64(samplePoints-1)*(t1-t0)
		y += 0.15 * (rng.NormFloat64() + 0.1*math.Sin(6*t))
		pts[i] = ui.ChartPoint{X: t, Y: y}
	}
	return pts
}

// dataSpikes: low-frequency signal with sharp positive pulses.
func dataSpikes() []ui.ChartPoint {
	return sampleSeries(0, 4*math.Pi, samplePoints, func(t float64) float64 {
		base := 0.3 * math.Sin(t)
		spike := 0.0
		for _, k := range []float64{1.0, 2.7, 3.6} {
			if math.Abs(t-k) < 0.12 {
				spike += 2.0 * math.Exp(-80*(t-k)*(t-k))
			}
		}
		return base + spike
	})
}

// dataDamped: exponentially decaying oscillation.
func dataDamped() []ui.ChartPoint {
	return sampleSeries(0, 4*math.Pi, samplePoints, func(t float64) float64 {
		return math.Exp(-0.1*t) * math.Sin(3*t)
	})
}

// dataSteps: piecewise level with a small HF ripple.
func dataSteps() []ui.ChartPoint {
	return sampleSeries(0, 10, samplePoints, func(t float64) float64 {
		return math.Mod(math.Floor(t), 5) + 0.1*math.Sin(20*t)
	})
}

func sampleSeries(t0, t1 float64, n int, f func(float64) float64) []ui.ChartPoint {
	pts := make([]ui.ChartPoint, n)
	for i := 0; i < n; i++ {
		t := t0
		if n > 1 {
			t = t0 + float64(i)/float64(n-1)*(t1-t0)
		}
		pts[i] = ui.ChartPoint{X: t, Y: f(t)}
	}
	return pts
}
