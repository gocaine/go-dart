package main

import (
	"fmt"
	"go-dart/cmd"
	"os"
)

func main() {
	fmt.Println("ready to Go !!")
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
