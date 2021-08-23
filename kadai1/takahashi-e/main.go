package main

import (
	"fmt"
	"gopherdojo-studyroom/kadai1/takahashi-e/mypkg"
	"os"
)
// main use
func main() {
	arguments, err := mypkg.ParseArguments()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if arguments.IsHelp {
		fmt.Println(mypkg.Help())
		os.Exit(0)
	}


	if err := arguments.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
