package main

import (
	"fmt"
	"os"

	"github.com/ugurkorkmaz/gogit"
)

const version string = "1.1.3"

func main() {
	var repo string
	var dir string

	gogit := gogit.New()
	if len(os.Args[1:]) > 0 {
		repo = os.Args[1]
	} else {
		fmt.Println("Usage: gogit <repo> [dir]")
		os.Exit(1)
	}
	if len(os.Args[2:]) > 0 {
		dir = os.Args[2]
	} else {
		dir = "."
	}

	if repo == "version" {
		fmt.Println("V" + version)
		os.Exit(0)
	}
	gogit.Run(repo, dir)
}
