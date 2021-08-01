package main

import (
	"flag"
	"fmt"
	"os"
	"tatsuo48/cmd"
)

type Options struct {
	src  string
	dest string
}

var options = Options{}

func init() {
	flag.StringVar(&options.src, "src", "jpg", "source file extension, default jpg")
	flag.StringVar(&options.dest, "dest", "png", "destination file extension, default png")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "specify a directory")
		os.Exit(1)
	}
	if err := cmd.Run(flag.Args()[0], options.src, options.dest); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
