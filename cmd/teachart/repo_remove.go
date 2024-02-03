// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
)

func NewRepoRemoveCmd(ctx context.Context, repoOptions *options.RepoOptions) *cobra.Command {
	repoRemoveOptions := options.NewRepoRemoveOptions(repoOptions)

	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: "Remove a chart repo.",
		Args:  ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRemoveOptions.Name = args[0]

			return runRepoRemove(ctx, repoOptions.Manager, repoRemoveOptions)
		},
	}

	return cmd
}

func runRepoRemove(ctx context.Context, manager *repo.Manager, opts *options.RepoRemoveOptions) error {
	return manager.Remove(opts.Name)
}
