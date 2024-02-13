// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/repo"
)

type repoRemoveOptions struct {
	*repoOptions

	name string
}

func NewRepoRemoveCmd(ctx context.Context, repoOpts *repoOptions) *cobra.Command {
	opts := &repoRemoveOptions{
		repoOptions: repoOpts,
	}

	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: "Remove a chart repo.",
		Args:  ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]

			return runRepoRemove(ctx, cmd, opts.manager, opts)
		},
	}

	return cmd
}

func runRepoRemove(ctx context.Context, cmd *cobra.Command, manager *repo.Manager, opts *repoRemoveOptions) error {
	err := manager.Remove(opts.name)
	if err == nil {
		fmt.Fprintf(cmd.OutOrStdout(), "repo %s removed\n", opts.name)
	}
	return err
}
