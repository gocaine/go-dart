package main

import (
	"fmt"
	"go-dart/cmd"
	"go-dart/server"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	server := server.NewServer()
	server.Start()

}
