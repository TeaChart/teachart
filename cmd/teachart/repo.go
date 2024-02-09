// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
)

type repoOptions struct {
	*options.GlobalOptions

	manager *repo.Manager
}

func NewRepoCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	opts := &repoOptions{
		GlobalOptions: globalOptions,
	}

	cmd := &cobra.Command{
		Use:   "repo add|remove|list",
		Short: "add, remove, list chart repos.",
		Args:  NoArgs,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			opts.manager = repo.NewManager(filepath.Join(globalOptions.GetInstallDir(), app.DefaultRepoDir), app.DefaultRemoteName)
			return opts.manager.Init()
		},
	}

	cmd.AddCommand(
		NewRepoAddCmd(ctx, opts),
		NewRepoRemoveCmd(ctx, opts),
		NewRepoListCmd(ctx, opts),
	)

	return cmd
}
