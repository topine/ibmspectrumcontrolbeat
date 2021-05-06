package main

import (
	"os"

	"github.com/topine/ibmspectrumcontrolbeat/cmd"

	_ "github.com/topine/ibmspectrumcontrolbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
