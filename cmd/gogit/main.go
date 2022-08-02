package main

import (
	"fmt"
	"os"

	"github.com/extendswork/gogit"
)

func main() {
	gogit := gogit.New()
	if len(os.Args[1:]) == 0 {
		fmt.Println("Usage: gogit <prefix> <dir>")
		os.Exit(1)
	}
	gogit.Run(os.Args[1], os.Args[2])
}
