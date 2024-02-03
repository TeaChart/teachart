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

func NewRepoAddCmd(ctx context.Context, repoOptions *options.RepoOptions) *cobra.Command {
	repoAddOptions := options.NewRepoAddOptions(repoOptions)

	cmd := &cobra.Command{
		Use:   "add NAME URL",
		Short: "Add a chart repo.",
		Args:  ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoAddOptions.Name = args[0]
			repoAddOptions.URL = args[1]

			return runRepoAdd(ctx, repoOptions.Manager, repoAddOptions)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&repoAddOptions.Force, "force", "f", false, "overwrite the repo if it already exists.")

	return cmd
}

func runRepoAdd(ctx context.Context, manager *repo.Manager, opts *options.RepoAddOptions) error {
	err := manager.Add(ctx, opts.Name, opts.URL, opts.Force)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		logrus.Errorf("repo `%s` already exists. use --force/-f to overwrite the repo.", opts.Name)
		return nil
	}
	return err
}
