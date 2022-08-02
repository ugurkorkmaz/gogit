package main

import (
	"fmt"
	"os"

	"github.com/extendswork/gogit"
)

func main() {
	var dir string
	gogit := gogit.New()
	if len(os.Args[1:]) == 0 {
		fmt.Println("Usage: gogit <prefix> <dir>")
		os.Exit(1)
	}
	if len(os.Args[2:]) == 0 {
		dir = "."
	}
	gogit.Run(os.Args[1], dir)
}
