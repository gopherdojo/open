# Gopher道場 自習室
[Gopher道場 自習室](https://gopherdojo.org/studyroom)用のリポジトリです。

## 課題の提出方法

1回目の課題を提出する場合は次のようにコードを書いて下さい。

* このリポジトリをフォークしてください
  * フォークしなくても直接プッシュできる場合はフォークせずに直接プッシュしてPRを送っても大丈夫です
* ブランチ名を`kadai1-tenntenn`のようにする
* `kadai1/tenntenn`のようにディレクトリを作る
* READMEに説明や文章による課題の回答を書く
* PRを送る

## レビューについて

* Slackの`#studyroom`チャンネルでレビューを呼びかけてみてください。
* PRは必ずレビューされる訳ではありません。
* 他の人のPRも積極的にレビューをしてみましょう！


## kadai1 課題の回答
* 【TRY】画像変換コマンドを作ろう

````
次の仕様を満たすコマンドを作って下さい
ディレクトリを指定する
指定したディレクトリ以下のJPGファイルをPNGに変換（デフォルト）
ディレクトリ以下は再帰的に処理する
変換前と変換後の画像形式を指定できる（オプション）
以下を満たすように開発してください
mainパッケージと分離する
自作パッケージと標準パッケージと準標準パッケージのみ使う
準標準パッケージ：golang.org/x以下のパッケージ
ユーザ定義型を作ってみる
GoDocを生成してみる
Go Modulesを使ってみる
````

kadai1ディレクトリ化で下記コマンドを実行

`````
go run . -s 変換するディレクトリ名 -f 変換前の拡張子 -cf 変換後のファイルタイプ
`````

- 使用例
````
go run main.go -s ./asset -f .png -cf .gif
````
を実行することで指定したディレクトリ（今回であれば./assset下に設備後に ```out.変換後の拡張子```がついた変換後の
ファイルが出力される。

* ヘルプ出力
```
go run main.go -help
```

* 自作パッケージ
* Go Modulesの使用　go mod init しています。
```aidl
. 
├── assets
│   ├── download-out.gif (変換後)
│   ├── download-out.jpg (変換後)
│   ├── download.jpg (変換前）
│   ├── download.png (変換前）
│   ├── neko\ test-out.gif (変換後)
│   ├── neko\ test-out.jpg (変換後)
│   ├── neko\ test.jpg (変換前）
│   └── neko\ test.png (変換前）
├── gopherdojo-studyroom
├── kadai1
├── main
├── main.go
└── mypkg
    ├── args.go (自作パッケージ)
    ├── convert.go (自作パッケ-ジ)
    └── mod
        └── cache
            └── lock


```

* ユーザー定義型として下記の構造体を作成しました。
```go
type Arguments struct {
	SelectedDirectory string
	SelectedFileType  string
	ConvertedFileType string
	StringPath        []string
	IsHelp            bool
	Args              []string
}
```

* GoDocの使用

Godocのインストール
````
go get golang.org/x/tools/cmd/godoc    
````

GoDocの起動
``
godoc -http=:6060 
``

``` http://localhost:6060 ```にアクセス

```mypkg```に記載されている。

以上