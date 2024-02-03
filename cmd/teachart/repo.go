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
	repoName string
}

func NewRepoCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	repoOptions := options.NewRepoOptions(globalOptions)

	cmd := &cobra.Command{
		Use:   "repo add|remove|list",
		Short: "add, remove, list chart repos.",
		Args:  NoArgs,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			repoOptions.Manager = repo.NewManager(filepath.Join(globalOptions.GetInstallDir(), app.DefaultRepoDir), app.DefaultRemoteName)
			return repoOptions.Manager.Init()
		},
	}

	cmd.AddCommand(
		NewRepoAddCmd(ctx, repoOptions),
		NewRepoRemoveCmd(ctx, repoOptions),
		NewRepoListCmd(ctx, repoOptions),
	)

	return cmd
}
