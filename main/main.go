package main

import (
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	//クラス1のプロトタイプ
	x1, y1 := 3.0, 6.0

	//クラス2のプロトタイプ
	x2, y2 := 8.0, 2.0

	//各クラスのサンプル
	n := 1000
	scatterData1 := randomPoints(n, x1, y1)
	scatterData2 := randomPoints(n, x2, y2)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//label
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData1)
	if err != nil {
		panic(err)
	}

	y, err := plotter.NewScatter(scatterData2)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	y.GlyphStyle.Color = color.RGBA{R: 155, B: 128, A: 255}
	p.Add(s)
	p.Add(y)
	p.Legend.Add("scatter", s)

	// Axis ranges that seem to include all bubbles.
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 0.5
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}

// randomPoints returns some random x, y points.
func randomPoints(n int, x, y float64) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {

		pts[i].X = random(x)

		pts[i].Y = random(y)
	}
	return pts
}
