# Golangでplot
## はじめに
golangでグラフを描く時，つまづいたので備忘録として残しておきます．
どうやらgonum plot`が便利らしい．
しかし，古い記事だと`go get`の場所が[google code](https://code.google.com/archive/p/plotinum/)を参照しているようだ．
現在は，`github`に移行しているので以下を実行．
```
go get gonum.org/v1/plot/...
```

## 二つのデータ群をplot
`(8.0, 2.0)`と`(3.0, 6.0)`を軸にそれぞれclass1,class2として100個,ガウス分布に従って，学習データを生成したものを`plot`
