# Pattern-recognition
## はじめに
現在パターン認識について学習している．　  
学習しているだけだとあれなので，実際にコードを書いてみることにした．  
機械学習なら`python`でしょ〜〜っていう世の中なので，今回は`golang`でやってみることにした．　　　　

## 今回行うこと．
まずは，ある学習パターンに対して，2クラスに分類するプログラムの実装を行う．

## 環境構築
`golang` の環境構築は割愛．  
`golang`でグラフを描くときには`gonum plot`が便利らしいので，以下でインストール
```
go get gonum.org/v1/plot/...
```

## 今回行うこと．
 `class1`と`class2`の二つの決定境界を決定する．

## 学習データについて
今回は，二次元の特徴空間について考える．まず，`(8.0, 2.0)`と`(3.0, 6.0)`を軸にそれぞれclass1,class2として，100個ガウス分布に従って，学習データを生成する．以下に，そのコードを示す．

```
//データ数
n := 100
//クラス1
class1 := randomPoints(n, 8.0, 2.0)
//クラス2
class2 := randomPoints(n, 3.0, 6.0)
//学習データの生成
func randomPoints(n int, x, y float64) (class [][]float64) {
	class := make([][]float64, n)
	for i := range class {
		class[i] = []float64{1.0, random(x), random(y)}
	}
	return
}
//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 1.0
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}
```

そして，生成された学習データを以下に示す．  
![sample](https://github.com/mytheta/Pattern-recognition/blob/master/main/first.png)

## 識別関数について
識別関数 $g(x) = W^tX$　のコードを以下に示す．

```
func innerProduct(w, x []float64) (f float64) {
	for i := range w {
		f += w[i] * x[i]
	}

	return
}
```
$g(x) = W^tX > 0$ならば$x\in\omega_1$  
$g(x) = W^tX < 0$ならば$x\in\omega_2$  
上の方法で，二つのクラスの決定境界を決める．コードを以下に示す．

```
func classify(w, x []float64) (t float64) {
	if innerProduct(w, x) >= 0 {
    //w1のクラス
		t = 1.0
	} else {
    //w2のクラス
		t = -1.0
	}
	return
}
```

## パーセプトロンの学習規則
パーセプトロンは以下の順序で，学習を行う．  
(1) 重みベクトルW の初期値を適当に選ぶ  
(2) 学習データの中から学習パターンを1つ選ぶ  
(3) 識別関数 g(x) = Wt X によって識別を行い,正しく識別できなかった場合のみ修正を行いW'を作る  
$$
W' = W ± \rho X (\rho > 0)
$$
(4) (2),(3)を学習データの全パターンで繰り返す  
(5) 学習データの全パターンを正しく認識すれば終了.  
誤りがあれば(2)へ戻る

 ---
まず，重みの初期値をきめる．

```
//重み
	w := []float64{20.0, -2.3, -2.1}
```
図を以下に示す.  
黄色の線が，初期の重みを表す．  
![sample](https://github.com/mytheta/Pattern-recognition/blob/master/main/weight.png)

以下のコードで，パーセプトロンの学習規則を表す．  
まず，class1を識別器にかけ，誤識別の場合，更新した重みと`-1`を返す．`fin`の変数で，`-1`を蓄積する．class2でも，同様に行う．  
もし，finの値が，`0`なら全て識別可能と判断し，ループから抜ける．`0`でなければ，`fin`を`0`にリセットし，再度，class1,class2を識別器にかける．
```
for {
  fin, count := 0, 0

  for i := 0; i < len(class1); i++ {
    w, count = train(w, class1[i], 1.0)
    fin += count
    fmt.Printf("w: %v\n", w)
  }

  for i := 0; i < len(class2); i++ {
    w, count = train(w, class2[i], -1.0)
    fin += count
    fmt.Printf("w: %v\n", w)
  }

  if fin == 0 {
    break
  }
}
```
(3)での誤識別の場合，以下のことを行う
$$
W' = W + \rho X (\omega_1のパターンを\omega_2と誤った時)  
$$
$$
W' = W - \rho X (\omega_2のパターンを\omega_1と誤った時)
$$
学習係数は，`p=0.7`とした．  

```
func train(w, x []float64, t float64) (nw []float64, fin int) {
	//学習係数
	p := 0.7
  //正解
	if classify(w, x) == t {
		fin = 0
		nw = w
    //誤識別
	} else {
		fin = -1
    //w1をw2と誤った場合
		if t == 1.0 {
      //w+px
			nw = add(w, multiple(p, x))
      //w2をw1と誤った場合
		} else {
			//w-px
			nw = subtraction(w, multiple(p, x))
		}
	}

	return
}
```
## 境界線のplot
パーセプトロンの収束定理を用いて更新した重みをplotする．
識別境界線を$g(X) = \omega_0 + \omega_1*x_1 + \omega_2*x_2 = 0$ で表し，以下にコードを示す．
```
lastBorder := plotter.NewFunction(func(x float64) float64 {
  //gx = wo + w1*x1 + w2*x2 = 0
  //x2 = -(w1 / w2)*x1 - w0 / w2
  return -(w[1]/w[2])*x - (w[0] / w[2])
})
```
線形分離可能なデータ群に対する識別境界線の図を以下に示す．  
![sample](https://github.com/mytheta/Pattern-recognition/blob/master/main/points.png)

## 重みの更新過程
重みの更新過程を以下に示す．

```
初期の重み: [20.0, -2.3, -2.1]
w: [19.3 -4.231054611569274 -5.97222820422448]
w: [20 1.9479066834312544 -4.586252176059295]
w: [19.3 0.016852071861980678 -8.458480380283774]
w: [20 5.256736901536651 -6.228807068574693]
最終の重み: [19.3 3.325682289967377 -10.101035272799173]
```

```
初期の重み: [20.0, -2.3, -2.1]

w: [19.3 -4.773782756117262 -5.596833430623128]
w: [20 1.8033769970206963 -4.706302335942745]
w: [19.3 -0.6704057590965657 -8.203135766565874]
w: [20 4.571607086639229 -6.214283837952144]
w: [19.3 2.0978243305219673 -9.711117268575272]
w: [20 8.625590924102898 -6.894440180967575]
w: [19.3 6.124313615142519 -12.005179129567711]
w: [18.6 3.7781696570763805 -14.307618506353363]
w: [19.3 10.305936250657313 -11.490941418745665]
w: [18.6 6.104553753036608 -15.6534278059303]
w: [19.3 9.979816187248032 -13.17262717550211]
w: [18.6 5.778433689627327 -17.335113562686743]
w: [19.3 11.499127522011715 -14.532639887292042]
w: [18.6 7.297745024391011 -18.695126274476674]
w: [19.3 11.173007458602434 -16.214325644048486]
w: [18.6 8.826863500536296 -18.516765020834136]
w: [19.3 13.788874472373523 -15.363364910025357]
w: [18.6 9.587491974752819 -19.52585129720999]
w: [19.3 14.549502946590046 -16.37245118640121]
w: [18.6 10.348120448969341 -20.534937573585843]
w: [19.3 15.310131420806568 -17.381537462777064]
最終の重み: [18.6 11.108748923185864 -21.544023849961697]

```
