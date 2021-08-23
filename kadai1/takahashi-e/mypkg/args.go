package mypkg

import (
	"flag"
)

//Argument Struct which use for command line flag.
type Arguments struct {
	SelectedDirectory string
	SelectedFileType  string
	ConvertedFileType string
	StringPath        []string
	IsHelp            bool
	Args              []string
}

//parse argument which provided by cli.
func ParseArguments() (*Arguments, error) {
	var argument Arguments

	flag.StringVar(&argument.SelectedDirectory, "s", "", "ディレクトリを指定")
	flag.StringVar(&argument.SelectedFileType, "f", ".jpg", "変換前のファイルタイプを指定")
	flag.StringVar(&argument.ConvertedFileType, "cf", ".png", "変換後のファイルタイプを指定")
	flag.BoolVar(&argument.IsHelp, "help", false, "display this help and exit")
	flag.Parse()
	return &argument, nil

}

//display options.
func Help() string {
	return `
Usage:
 convert [options] command
	Options:
  	-s,  変換したいファイルがあるディレクトリを指定
  	-f,  変換前のファイルタイプを.jpeg , .png , .gifを指定 デフォルトは .jpg 
  	-cf, 変換後のファイルタイプを .png , .gif 指定 デフォルトは.png`
}
