// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/compose-spec/compose-go/v2/cli"
	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/compose"
	"github.com/yp05327/teachart/pkg/engine"
	"github.com/yp05327/teachart/pkg/options"
	"helm.sh/helm/v3/pkg/cli/values"
)

type lintOptions struct {
	*options.GlobalOptions
	*options.ChartOptions
	*configOptions

	values    values.Options
	chartDirs []string
}

type configOptions struct {
	Format              string
	resolveImageDigests bool
	noInterpolate       bool
	noNormalize         bool
	noResolvePath       bool
	noConsistency       bool
}

func NewLintCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	opts := lintOptions{
		GlobalOptions: globalOptions,
		ChartOptions:  options.NewChartOptions(globalOptions.GetRepoRootDir()),
		configOptions: &configOptions{},
		values:        values.Options{},
		chartDirs:     []string{"."},
	}

	cmd := &cobra.Command{
		Use:   "lint [PATHS]",
		Short: "Lint charts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.chartDirs = args
			}
			return runLint(ctx, cmd, opts)
		},
	}

	flags := cmd.Flags()
	addChartFlags(flags, opts.ChartOptions)
	addValueFlags(flags, &opts.values)
	addConfigFlasg(flags, opts.configOptions)

	return cmd
}

func runLint(ctx context.Context, cmd *cobra.Command, opts lintOptions) error {
	pofs := []cli.ProjectOptionsFn{
		cli.WithInterpolation(!opts.noInterpolate),
		cli.WithResolvedPaths(!opts.noResolvePath),
		cli.WithNormalization(!opts.noNormalize),
		cli.WithConsistency(!opts.noConsistency),
		cli.WithDiscardEnvFile,
		cli.WithContext(ctx),
	}

	lintFn := func(chartDir string) ([]byte, error) {
		// create a temp dir to save rendered files
		tempDir, err := os.MkdirTemp("", "teachart-lint")
		if err != nil {
			return nil, errors.Wrap(err, "Create temp directory error")
		}
		if !opts.Debug {
			defer os.RemoveAll(tempDir)
		}

		logrus.Debugf("Temp directory created:%s", tempDir)

		// render templates
		renderEngine, err := engine.NewRenderEngine(chartDir, opts.GetTeaChart(), true)
		if err != nil {
			return nil, errors.Wrap(err, "Create helm engine error")
		}
		files, err := renderEngine.Render(opts.values, true)
		if err != nil {
			return nil, errors.Wrap(err, "Render chart error")
		}
		// run compose config
		client, err := compose.NewClient(nil, &compose_cmd.ProjectOptions{
			ProjectName: "teachart-lint",
			ConfigPaths: renderEngine.GetConfigPaths(files),
		})
		if err != nil {
			return nil, errors.Wrap(err, "Create compose client error")
		}
		return client.Config(ctx, pofs, api.ConfigOptions{
			Format:              opts.Format,
			ResolveImageDigests: opts.resolveImageDigests,
		})
	}

	failed := 0
	for _, chartDir := range opts.chartDirs {
		fmt.Fprintf(cmd.OutOrStdout(), "==> Linting %s\n", chartDir)

		content, err := lintFn(chartDir)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error %s\n", err)
			failed++
			continue
		}

		if !opts.Quiet {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(content))
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%d chart(s) linted, %d chart(s) failed\n", len(opts.chartDirs), failed)
	return nil
}
