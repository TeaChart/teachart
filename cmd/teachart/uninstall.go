// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"
	"os"

	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/compose"
	"github.com/yp05327/teachart/pkg/options"
)

type uninstallOptions struct {
	*options.GlobalOptions
	*compose_cmd.ProjectOptions
	*downOptions

	cleanTemp bool
}

type downOptions struct {
	removeOrphans bool
	timeChanged   bool
	timeout       int
	volumes       bool
	images        string
}

func NewUninstallCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	opts := &uninstallOptions{
		GlobalOptions: globalOptions,
		downOptions:   &downOptions{},
	}

	cmd := &cobra.Command{
		Use:        "uninstall",
		Aliases:    []string{"del", "delete", "un"},
		SuggestFor: []string{"remove", "rm"},
		Short:      "uninstall a deploy",
		Args:       NoArgs,
		PreRunE: compose_cmd.AdaptCmd(func(ctx context.Context, cmd *cobra.Command, args []string) error {
			opts.timeChanged = cmd.Flags().Changed("timeout")
			if opts.images != "" {
				if opts.images != "all" && opts.images != "local" {
					return fmt.Errorf("invalid value for --rmi: %q", opts.images)
				}
			}
			return nil
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: support apply services
			services := []string{}
			client, err := compose.NewClient(services, opts.GetProjectOptions())
			if err != nil {
				return errors.Wrap(err, "Create compose client error")
			}

			return runUninstall(ctx, client, opts)
		},
	}

	flags := cmd.Flags()
	addDownFlags(flags, opts.downOptions)
	flags.BoolVar(&opts.cleanTemp, "clean-temp", true, "cleanup the temp folder after remove all containers. default is true")

	return cmd
}

func runUninstall(ctx context.Context, client *compose.Client, opts *uninstallOptions) error {
	down := api.DownOptions{}
	if err := client.Down(ctx, down); err != nil {
		return errors.Wrap(err, "Uninstall error")
	}

	// clean temp folder
	if opts.cleanTemp {
		if err := os.RemoveAll(opts.GetTempDir()); err != nil {
			return errors.Wrap(err, "Remove temp directory error")
		}
	}
	return nil
}
