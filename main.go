package main

import (
	"fmt"
	"github.com/gocaine/go-dart/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(-1)
	}
}
