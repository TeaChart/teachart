// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
)

type repoAddOptions struct {
	*repoOptions
	*options.ChartOptions

	force bool
}

func NewRepoAddCmd(ctx context.Context, repoOpts *repoOptions) *cobra.Command {
	opts := &repoAddOptions{
		repoOptions:  repoOpts,
		ChartOptions: options.NewChartOptions(repoOpts.GetRepoRootDir()),
	}

	cmd := &cobra.Command{
		Use:   "add NAME URL",
		Short: "Add a chart repo.",
		Args:  ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.URL = args[1]

			return runRepoAdd(ctx, opts.manager, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.force, "force", "f", false, "overwrite the repo if it already exists.")

	return cmd
}

func runRepoAdd(ctx context.Context, manager *repo.Manager, opts *repoAddOptions) error {
	err := manager.Add(ctx, opts.Name, opts.URL, opts.force)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		logrus.Errorf("repo `%s` already exists. use --force/-f to overwrite the repo.", opts.Name)
		return nil
	}
	return err
}
