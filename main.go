package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/ics-availability-server/cmd"
	"github.com/gabe565/ics-availability-server/internal/config"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)
	root := cmd.New(cmd.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
