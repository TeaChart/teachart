// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"

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

			return runRepoRemove(ctx, opts.manager, opts)
		},
	}

	return cmd
}

func runRepoRemove(ctx context.Context, manager *repo.Manager, opts *repoRemoveOptions) error {
	return manager.Remove(opts.name)
}
