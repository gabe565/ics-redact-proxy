package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "IAS_"

var ErrNoSource = errors.New("no source URL defined")

func (c *Config) Load(cmd *cobra.Command) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	initLog(cmd)
	if err := errors.Join(errs...); err != nil {
		return err
	}

	if c.SourceURL == "" {
		return ErrNoSource
	}
	return nil
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
