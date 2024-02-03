// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
)

func NewRepoListCmd(ctx context.Context, repoOptions *options.RepoOptions) *cobra.Command {
	repoListOptions := options.NewRepoListOptions(repoOptions)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all chart repos.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRepoList(ctx, repoOptions.Manager, repoListOptions)
		},
	}

	return cmd
}

func runRepoList(ctx context.Context, manager *repo.Manager, opts *options.RepoListOptions) error {
	reposMap, err := manager.List()
	if err != nil {
		return err
	}

	if len(reposMap) == 0 {
		fmt.Println("no repositories. use `repo add` to add repos")
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
	fmt.Println(table.String())
	return nil
}
