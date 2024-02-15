// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

var chartName = regexp.MustCompile("^[a-zA-Z0-9._-]+$")

const maxChartNameLength = 250

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
		Use:   "init",
		Short: "Initialize a teachart",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(initOptions.ChartOptions.Name) == 0 {
				initOptions.ChartOptions.Name = globalOptions.GetProjectName()
			}
			return runInit(ctx, cmd, initOptions)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&initOptions.Name, "name", "n", "", "teachart name.")

	return cmd
}

func runInit(ctx context.Context, cmd *cobra.Command, opts *initOptions) error {
	// Sanity-check the name of a chart so user doesn't create one that causes problems.
	if opts.ChartOptions.Name == "" || len(opts.ChartOptions.Name) > maxChartNameLength {
		return fmt.Errorf("chart name must be between 1 and %d characters", maxChartNameLength)
	}
	if !chartName.MatchString(opts.Name) {
		return fmt.Errorf("chart name must match the regular expression %q", chartName.String())
	}

	path, err := filepath.Abs(opts.GetProjectDir())
	if err != nil {
		return err
	}

	if fi, err := os.Stat(path); err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.Errorf("no such directory %s", path)
	}

	cdir := filepath.Join(path, opts.ChartOptions.Name)
	if fi, err := os.Stat(cdir); err == nil && !fi.IsDir() {
		return errors.Errorf("file %s already exists and is not a directory", cdir)
	}

	files := []struct {
		path    string
		content string
	}{
		{
			// Chart.yaml
			path:    filepath.Join(cdir, app.DefaultMetadataFileName),
			content: fmt.Sprintf(app.DefaultMetadataFile, opts.Name),
		},
		{
			// values.yaml
			path:    filepath.Join(cdir, app.DefaultValuesFileName),
			content: fmt.Sprintf(app.DefaultValuesFile, opts.Name),
		},
		{
			// NOTES.txt
			path:    filepath.Join(cdir, app.DefaultNotesFileName),
			content: fmt.Sprintf(app.DefaultNotesFile, opts.Name),
		},
	}

	for _, file := range files {
		if _, err := os.Stat(file.path); err == nil {
			// There is no handle to a preferred output stream here.
			fmt.Fprintf(cmd.OutOrStdout(), "WARNING: File %q already exists. Overwriting.\n", file.path)
		}
		if err := os.MkdirAll(filepath.Dir(file.path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			return err
		}
	}
	// add the TemplatesDir
	if err := os.MkdirAll(filepath.Join(cdir, app.DefaultTemplatesDir), 0755); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s created\n", opts.Name)
	return nil
}
