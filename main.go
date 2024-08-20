package main

import (
	"os"

	"github.com/gabe565/ics-availability-server/cmd"
)

var version = "beta"

func main() {
	root := cmd.New(cmd.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
