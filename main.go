package main

import (
	"os"

	"github.com/gabe565/ics-availability-server/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
