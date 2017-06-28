package main

import (
	"os"

	"github.com/ashmckenzie/go-littlefly/littlefly"
	"github.com/ashmckenzie/go-littlefly/littlefly/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		littlefly.Log.Fatalf("%s\n", err)
		os.Exit(1)
	}
}
