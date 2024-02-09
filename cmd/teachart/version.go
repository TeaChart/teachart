// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/options"
	"go.szostok.io/version"
)

func NewVersionCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	versionOptions := options.NewVersionOptions(globalOptions)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version information",
		Args:  NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(ctx, versionOptions)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&versionOptions.Short, "short", false, "only print version number")
	flags.StringVar(&versionOptions.Format, "format", "text", "format for version string, can be text/json/yaml. default is `text`")

	return cmd
}

func runVersion(ctx context.Context, opts *options.VersionOptions) error {
	versionInfo := version.Get()

	switch opts.Format {
	case "json":
		json, err := versionInfo.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, json)
	case "yaml":
		yaml, err := versionInfo.MarshalYAML()
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, yaml)
	default:
		if opts.Short {
			fmt.Println(versionInfo.Version)
		} else {
			fmt.Printf("%#v\n", versionInfo)
		}
	}
	return nil
}
