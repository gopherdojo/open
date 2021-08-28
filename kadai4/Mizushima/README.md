## 課題4 おみくじAPIを作ってみよう
- JSON形式でおみくじの結果を返す
- 正月（1/1-1/3）だけ大吉にする
- ハンドラのテストを書いてみる


### インストール方法
```bash
go get github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai4/Mizushima
```

### 使い方
1. 実行ファイル作成 & テスト
```bash
$ make
```
2. サーバーの起動
```bash
$ ./bin/omikuji
```
3. おみくじ結果の表示
 ブラウザから"http:localhost:8080"にアクセスすると、下記のような表示が出ます。  
 例：{"result":"大吉"}
  

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
│   └── omikuji # 'make' command required.
├── fortune_server.go
├── fortune_server_test.go
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── README.md
```