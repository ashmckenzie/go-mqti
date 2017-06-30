package main

import (
	"os"

	"github.com/ashmckenzie/go-mqti/mqti"
	"github.com/ashmckenzie/go-mqti/mqti/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		mqti.Log.Fatalf("%s\n", err)
		os.Exit(1)
	}
}
