package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/gabe565/ics-availability-server/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func New(opts ...Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ics-availability-server",
		Short: "Fetches an ics file and redacts all data except for configured fields.",
		RunE:  run,

		DisableAutoGenTag: true,
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

	log.Info().
		Str("version", cmd.Annotations[VersionKey]).
		Str("commit", cmd.Annotations[CommitKey]).
		Msg("ICS availability server")

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	return server.ListenAndServe(ctx, conf)
}
