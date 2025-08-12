package main

import (
	"log/slog"
	"os"

	"gabe565.com/ics-redact-proxy/cmd"
	"gabe565.com/ics-redact-proxy/internal/config"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/slogx"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
