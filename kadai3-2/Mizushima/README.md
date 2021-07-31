## 課題3-2 分割ダウンローダを作ろう
- 分割ダウンロードを行う
- Rangeアクセスを用いる
- いくつかのゴルーチンでダウンロードしてマージする
- エラー処理を工夫する
  - golang.org/x/sync/errgourpパッケージなどを使ってみる
- キャンセルが発生した場合の実装を行う



### コマンドラインオプション

 | ショートオプション | ロングオプション | 説明 | デフォルト |
 | --------- | --------- | --------- | --------- |
 | -h | --help | 使い方を表示して終了 | - |
 | -p \<num> | --procs \<num> | プロセス数を指定 | お使いのPCのコア数 |
 | -o \<path> | --output \<path> | ダウンロードしたファイルをどこのディレクトリに保存するか指定する | カレントディレクトリ |
 | -t \<num> | --timeout \<num> | サーバーへのリクエストを止める時間を秒数で指定 | 120 |


### インストール方法
```bash
go get github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima
```

### 使い方
1. 実行ファイル作成
```bash
$ make build
```
2. URLを指定してダウンロード実行
```bash
$ ./bin/paraDW [option] URL( URL URL ...)
```

※ URLは複数指定できます  
※ ダウンロード先がRangeアクセスに対応していれば、go routineを使った並行ダウンロードを行い、そうでなければ１プロセスのダウンロードを行います  

### テストの方法
- バイナリビルド & テスト
```bash
$ make
```
- テスト後の処理(掃除)
```bash
$ make clean
```

### ディレクトリ構成
```bash
.
├── bin  
│   └── paraDW  # "make build" command required.  
├── download  
│   ├── download.go  
│   ├── download_test.go  
│   ├── go.mod  
│   └── go.sum  
├── getheader  
│   ├── geHeader_test.go  
│   ├── getHeader.go  
│   └── go.mod  
├── .gitignore  
├── go.mod  
├── go.sum  
├── listen  
│   ├── listen.go  
│   └── listen_test.go  
├── main.go  
├── Makefile  
├── README.md  
├── request  
│   ├── go.mod  
│   ├── request.go  
│   └── request_test.go  
└── testdata  
    ├── 003  
    └── z4d4kWk.jpg  
```

### 参考にしたもの
[pget](https://qiita.com/codehex/items/d0a500ac387d39a34401)  (goroutineを使ったダウンロード処理、コマンドラインオプションの処理等々)  
https://github.com/gopherdojo/dojo3/pull/50  (ctrl+cを押したときのキャンセル処理など)