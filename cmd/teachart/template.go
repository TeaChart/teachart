// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/engine"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
	"helm.sh/helm/v3/pkg/cli/values"
)

type templateOptions struct {
	*options.GlobalOptions
	*options.ChartOptions

	save   bool
	values values.Options
	debug  bool
}

func NewTemplateCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	opts := templateOptions{
		GlobalOptions: globalOptions,
		ChartOptions:  &options.ChartOptions{},
		values:        values.Options{},
		debug:         globalOptions.Debug,
	}

	cmd := &cobra.Command{
		Use:   "template REPO_NAME",
		Short: "Render templates locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			} else if opts.Dir == "" {
				return fmt.Errorf("should provide a repo name or a dir path to the teachart repo\n\nUsage:%s", cmd.UseLine())
			}

			if opts.Dir == "" {
				manager := repo.NewManager(filepath.Join(globalOptions.GetInstallDir(), app.DefaultRepoDir), app.DefaultRemoteName)
				if err := manager.Init(); err != nil {
					return err
				}
				// prepare repo
				authMethod, err := opts.GetAuthMethod()
				if err != nil {
					return err
				}
				if err := manager.Checkout(ctx, opts.Name, authMethod, repo.CheckOutOptions{
					Version: opts.Version,
					Commit:  opts.Commit,
					Debug:   opts.debug,
				}); err != nil {
					return errors.Wrap(err, "checkout failed")
				}
			}
			return runTemplate(ctx, opts)
		},
	}

	flags := cmd.Flags()
	addValueFlags(flags, &opts.values)
	addChartFlags(flags, opts.ChartOptions)
	flags.BoolVar(&opts.save, "save", false, "save rendered files into the temp dir. default is false")

	return cmd
}

func runTemplate(ctx context.Context, opts templateOptions) error {
	renderEngine, err := engine.NewRenderEngine(opts.GetChartDir(), opts.GetTeaChart(), nil)
	if err != nil {
		return errors.Wrap(err, "Create helm engine error")
	}
	_, err = renderEngine.Render(opts.values, opts.save)
	return err
}
