package main

import (
	"os"
	"path/filepath"
	"fmt"
	"flag"
	"imgconv/imgconv"
)

func main() {
	var beforeExt = flag.String("before", "jpg", "Extension before conversion")
	var afterExt = flag.String("after", "png", "Extension after conversion")
	flag.Parse()
	dirs := flag.Args()

	if *beforeExt == *afterExt {
		fmt.Printf("beforeとafterに同じ拡張子%sが指定されています", *beforeExt)
		return
	}

	for _, dir := range dirs {
		dirwalk(dir, *beforeExt, *afterExt)
	}
}

// ディレクトリの中を巡回する
func dirwalk(dir, beforeExt, afterExt string){
	// ディレクトリの存在をチェック
	if _, err := os.Stat(dir); err != nil {
    fmt.Printf("%v\n", err)
		return
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error{
		// beforeExtと一致する時だけ変換を行う
		if filepath.Ext(path) == "." + beforeExt {
			fmt.Println(path)
			fmt.Println(dir)
			imgconv.Imgconv(path, afterExt)
		}

		return nil
	})
}

