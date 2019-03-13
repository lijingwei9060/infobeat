package main

import (
	"os"

	"github.com/lijingwei9060/infobeat/cmd"

	_ "github.com/lijingwei9060/infobeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
