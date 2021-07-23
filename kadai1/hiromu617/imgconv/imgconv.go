package imgconv

import (
	"image"
	"regexp"
	"image/jpeg"
	"image/png"
	"image/gif"
	"fmt"
	"os"
	"path/filepath"
)

func Imgconv(path, afterExt string){
	// fileを開く
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer f.Close()

	// 元画像をデコード
	img, _, err := image.Decode(f) 
	if err != nil {
		fmt.Println("decode:", err)
		return
	}

	// 拡張子を除くファイル名を取得
	rep := regexp.MustCompile(filepath.Ext(path) + "$")
  fileName := filepath.Base(rep.ReplaceAllString(path, ""))

	// 変換先のファイルを作成
	fso, err := os.Create(filepath.Join(filepath.Dir(path), fileName + "." + afterExt))
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()
	
	// 拡張子によって、それぞれの形式に変換
	switch afterExt {
	case "jpeg", "jpg":
		jpeg.Encode(fso, img, &jpeg.Options{})
	case "png":
		png.Encode(fso, img)
	case "gif":
		gif.Encode(fso, img, nil)
	}
}