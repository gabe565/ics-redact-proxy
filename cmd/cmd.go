package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/gabe565/ics-availability-server/internal/server"
	"github.com/spf13/cobra"
)

func New(opts ...Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ics-availability-server",
		Short: "Fetches an ics file and redacts all data except for configured fields.",
		RunE:  run,

		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd.Flags())
	cmd.SetContext(config.NewContext(context.Background(), conf))

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("command missing config")
	}

	if err := conf.Load(cmd); err != nil {
		return err
	}

	slog.Info("ICS availability server",
		"version", cmd.Annotations[VersionKey],
		"commit", cmd.Annotations[CommitKey],
	)

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	return server.ListenAndServe(ctx, conf)
}
