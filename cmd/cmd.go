package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/ics-redact-proxy/internal/config"
	"gabe565.com/ics-redact-proxy/internal/server"
	"gabe565.com/utils/cobrax"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ics-redact-proxy",
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

	slog.Info("ICS redact proxy", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))

	if !conf.NoVerify {
		if err := verifySource(cmd.Context(), conf); err != nil {
			return err
		}
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	return server.ListenAndServe(ctx, conf)
}

var ErrSourceVerify = errors.New("source verification failed")

func verifySource(ctx context.Context, conf *config.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, conf.SourceURL, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSourceVerify, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSourceVerify, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrSourceVerify, resp.Status)
	}

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSourceVerify, err)
	}

	//nolint:gosec
	slog.Info("Upstream verification succeeded", "status", resp.Status, "size", humanize.IBytes(uint64(n)))
	return nil
}
