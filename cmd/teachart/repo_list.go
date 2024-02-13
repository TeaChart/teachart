// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
)

type repoListOptions struct {
	*repoOptions
}

func NewRepoListCmd(ctx context.Context, repoOpts *repoOptions) *cobra.Command {
	opts := &repoListOptions{
		repoOptions: repoOpts,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all chart repos.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRepoList(ctx, cmd, opts)
		},
	}

	return cmd
}

func runRepoList(ctx context.Context, cmd *cobra.Command, opts *repoListOptions) error {
	reposMap, err := opts.manager.List()
	if err != nil {
		return err
	}

	if len(reposMap) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no repositories. use `repo add` to add repos")
		return nil
	}

	table := uitable.New()
	table.AddRow("NAME", "URL")
	for name, repo := range reposMap {
		remote, err := repo.Remote(app.DefaultRemoteName)
		if err != nil {
			return err
		}
		table.AddRow(name, remote.Config().URLs[0])
	}
	fmt.Fprintln(cmd.OutOrStdout(), table.String())
	return nil
}
