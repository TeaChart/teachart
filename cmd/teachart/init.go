// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/options"
	"helm.sh/helm/v3/pkg/chartutil"
)

type initOptions struct {
	*options.GlobalOptions
	*options.ChartOptions
}

func NewInitCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	initOptions := &initOptions{
		GlobalOptions: globalOptions,
		ChartOptions:  &options.ChartOptions{},
	}

	cmd := &cobra.Command{
		Use:   "init [NAME]",
		Short: "Initialize the teachart. Default name is folder name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				initOptions.Name = args[0]
			}

			return runInit(ctx, initOptions)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&initOptions.Name, "name", "n", "", "teachart name.")

	return cmd
}

func runInit(ctx context.Context, opts *initOptions) error {
	// TODO

	dir, err := chartutil.Create(opts.Name, opts.GetProjectDir())
	fmt.Println(dir)
	return err
}
