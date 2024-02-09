// Copyright © 2024 TeaChart Authors

package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

var rootDesc = "A tool for deploying multiple containers in one shot."

func NewRootCmd(ctx context.Context, globalOptions *options.GlobalOptions) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "teachart",
		Short:         rootDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			logLevel, err := logrus.ParseLevel(globalOptions.LogLevel)
			if err != nil {
				return err
			}

			switch {
			case globalOptions.Debug:
				logLevel = logrus.TraceLevel
			case globalOptions.Quiet:
				logLevel = logrus.WarnLevel
			}
			logrus.SetLevel(logLevel)
			return nil
		},
	}
	cmd.AddCommand(
		NewInitCmd(ctx, globalOptions),
		NewTemplateCmd(ctx, globalOptions),
		NewLintCmd(ctx, globalOptions),

		NewRepoCmd(ctx, globalOptions),
		NewInstallCmd(ctx, globalOptions),
		NewUninstallCmd(ctx, globalOptions),

		NewVersionCmd(ctx, globalOptions),
	)

	flags := cmd.PersistentFlags()
	flags.StringVar(&globalOptions.ComposeExecutable, "compose-bin", app.DefaultComposeExecutable, "path to the docker compose executable file.")
	flags.BoolVarP(&globalOptions.Quiet, "quiet", "q", false, "Set log-level to warn.")
	flags.BoolVar(&globalOptions.Debug, "debug", false, "Set log-level to debug, this disables --quiet/-q.")
	flags.StringVar(&globalOptions.LogLevel, "log-level", "info", "Set log level, default is \"info\". Will be ignored when --debug or --quiet/-q is set.")
	flags.StringVar(&globalOptions.ProjectDir, "project-directory", "", "project directory for the docker compose.")
	flags.StringVarP(&globalOptions.ProjectName, "project-name", "p", "", "path to the docker compose executable file.")
	flags.StringVar(&globalOptions.TempDirName, "temp-dir-name", app.DefaultTempDirName, "temp path for the generated docker compose files.")
	flags.BoolP("help", "h", false, "help for teachart")

	flags.ParseErrorsWhitelist.UnknownFlags = true

	return cmd, nil
}
