package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	//重み
	w := []float64{20.0, -2.3, -2.1}

	// 図の生成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//任意の点
	dots := make(plotter.XYs, 2)

	//クラス1
	x1, y1 := 8.0, 2.0
	dots[0].X = x1
	dots[0].Y = y1

	//クラス2
	x2, y2 := 3.0, 6.0
	dots[1].X = x2
	dots[1].Y = y2

	//各クラスのサンプル
	n := 200
	class1, plotdata1 := randomPoints(n, x1, y1)
	class2, plotdata2 := randomPoints(n, x2, y2)

	//初期重みのplot----------------------------------------------
	firstW := w

	border := plotter.NewFunction(func(x float64) float64 {
		return -(firstW[1]/firstW[2])*x - (firstW[0] / firstW[2])
	})
	border.Color = color.RGBA{B: 255, A: 45}
	p.Add(border)
	//----------------------------------------------------------

	times := 0
	//重みの更新が終わるまで学習----------------------------------------------------------
	for {
		fin, count := 0, 0

		for i := 0; i < len(class2); i++ {
			w, count = train(w, class2[i], -1.0)
			fin += count
			//	fmt.Printf("w: %v\n", w)
		}
		//学習
		for i := 0; i < len(class1); i++ {
			w, count = train(w, class1[i], 1.0)
			fin += count
			// fmt.Printf("w: %v\n", w)
		}

		times++
		if fin == 0 {
			break
		}
	}
	//------------------------------------------------------------------------------

	//最終境界線のplot--------------------------------------------------------
	lastBorder := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(w[1]/w[2])*x - (w[0] / w[2])
	})
	lastBorder.Color = color.RGBA{B: 255, A: 255}
	//----------------------------------------------------------------------

	//label
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(plotdata1)
	if err != nil {
		panic(err)
	}

	y, err := plotter.NewScatter(plotdata2)
	if err != nil {
		panic(err)
	}

	r, err := plotter.NewScatter(dots)
	if err != nil {
		panic(err)
	}
	fmt.Printf("学習終わり")
	fmt.Printf("times: %d\n", times)
	fmt.Printf("w: %v\n", w)

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 55}
	y.GlyphStyle.Color = color.RGBA{R: 155, B: 128, A: 255}
	r.GlyphStyle.Color = color.RGBA{R: 128, B: 0, A: 0}
	p.Add(s)
	p.Add(y)
	p.Add(r)
	p.Add(lastBorder)
	p.Legend.Add("class1", s)
	p.Legend.Add("class2", y)

	// Axis ranges
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "report.png"); err != nil {
		panic(err)
	}

}

//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 1.0
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}

//学習データの生成
func randomPoints(n int, x, y float64) ([][]float64, plotter.XYs) {
	matrix := make([][]float64, n)
	pts := make(plotter.XYs, n)
	for i := range matrix {
		l := random(x)
		m := random(y)
		matrix[i] = []float64{1.0, l, m}
		pts[i].X = l
		pts[i].Y = m
	}
	return matrix, pts
}

//学習
func train(w, x []float64, t float64) (nw []float64, fin int) {
	//学習係数
	p := 0.7
	if classify(w, x) == t {
		fin = 0
		nw = w
	} else {
		fin = 1
		//w1をw2と誤った場合->w=w+px(t=1.0)
		//w2をw1と誤った場合->w=w-px(t=-1.0)
		p = t * p
		nw = add(w, multiple(p, x))

		//更新した重みをplot
		fmt.Printf("w: %v\n", nw)
	}

	return
}

//識別関数
func classify(w, x []float64) (t float64) {
	if innerProduct(w, x) >= 0 {
		//w1
		t = 1.0
	} else {
		//w2
		t = -1.0
	}
	return
}

//識別計算
func innerProduct(w, x []float64) (f float64) {
	if len(w) != len(x) {
		panic("エラーですよ")
	}

	for i := range w {
		f += w[i] * x[i]
	}

	return
}

//w+px
func add(w, x []float64) (z []float64) {
	if len(w) != len(x) {
		panic("エラーですよ")
	}

	for i, _ := range w {
		z = append(z, w[i]+x[i])
	}
	return
}

//±p*x
func multiple(p float64, x []float64) (y []float64) {
	for i := range x {
		y = append(y, p*x[i])
	}
	return
}
