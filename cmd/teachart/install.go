// Copyright © 2024 TeaChart Authors

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/compose"
	"github.com/yp05327/teachart/pkg/engine"
	"github.com/yp05327/teachart/pkg/options"
	"github.com/yp05327/teachart/pkg/repo"
	"helm.sh/helm/v3/pkg/cli/values"
)

type installOptions struct {
	*options.GlobalOptions
	*options.ChartOptions
	*createOptions
	*startOptions

	values values.Options
}

type createOptions struct {
	Pull          string
	pullChanged   bool
	removeOrphans bool
	ignoreOrphans bool
	forceRecreate bool
	noRecreate    bool
	recreateDeps  bool
	noInherit     bool
	timeChanged   bool
	timeout       int
	quietPull     bool
	scale         []string
}

type startOptions struct {
	noStart     bool
	noDeps      bool
	noColor     bool
	noPrefix    bool
	timestamp   bool
	wait        bool
	waitTimeout int
}

func (c *createOptions) recreateStrategy() string {
	if c.noRecreate {
		return api.RecreateNever
	}
	if c.forceRecreate {
		return api.RecreateForce
	}
	return api.RecreateDiverged
}

func (c *createOptions) dependenciesRecreateStrategy() string {
	if c.noRecreate {
		return api.RecreateNever
	}
	if c.recreateDeps {
		return api.RecreateForce
	}
	return api.RecreateDiverged
}

func (c *createOptions) getTimeOut() *time.Duration {
	if c.timeChanged {
		t := time.Duration(c.timeout) * time.Second
		return &t
	}
	return nil
}

func (s *startOptions) GetWaitTimeOut() time.Duration {
	return time.Duration(s.waitTimeout) * time.Second
}

func NewInstallCmd(ctx context.Context, globalOptions *options.GlobalOptions) *cobra.Command {
	create := &createOptions{}
	start := &startOptions{}
	opts := &installOptions{
		GlobalOptions: globalOptions,
		ChartOptions:  options.NewChartOptions(globalOptions.GetRepoRootDir()),
		createOptions: create,
		startOptions:  start,
		values:        values.Options{},
	}

	cmd := &cobra.Command{
		Use:   "install REPO_NAME",
		Short: "Install a chart",
		PreRunE: compose_cmd.AdaptCmd(func(ctx context.Context, cmd *cobra.Command, args []string) error {
			create.pullChanged = cmd.Flags().Changed("pull")
			create.timeChanged = cmd.Flags().Changed("timeout")
			return validateFlags(opts)
		}),
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
					Debug:   opts.Debug,
				}); err != nil {
					return errors.Wrap(err, "Repo checkout failed")
				}
			}
			// TODO: support apply services
			services := []string{}

			client, err := compose.NewClient(services, opts.GetProjectOptions())
			if err != nil {
				return errors.Wrap(err, "Create compose client error")
			}
			renderEngine, err := engine.NewRenderEngine(opts.GetChartDir(), opts.GetTeaChart())
			if err != nil {
				return errors.Wrap(err, "Create helm engine error")
			}
			return runInstall(ctx, client, renderEngine, opts)
		},
	}

	flags := cmd.Flags()
	addChartFlags(flags, opts.ChartOptions)
	addValueFlags(flags, &opts.values)
	addCreateFlags(flags, opts.createOptions)
	addStartFlags(flags, opts.startOptions)

	return cmd
}

func validateFlags(opts *installOptions) error {
	if opts.forceRecreate && opts.noRecreate {
		return fmt.Errorf("--force-recreate and --no-recreate are incompatible")
	}
	if opts.recreateDeps && opts.noRecreate {
		return fmt.Errorf("--always-recreate-deps and --no-recreate are incompatible")
	}
	return nil
}

func runInstall(ctx context.Context, client *compose.Client, renderEngine engine.RenderEngine, opts *installOptions) error {
	// render templates
	files, err := renderEngine.Render(opts.values, true)
	if err != nil {
		return errors.Wrap(err, "Render chart error")
	}

	// set config paths
	client.ConfigPaths = renderEngine.GetConfigPaths(files)

	// run compose up
	create := api.CreateOptions{
		RemoveOrphans:        opts.removeOrphans,
		IgnoreOrphans:        opts.ignoreOrphans,
		Recreate:             opts.recreateStrategy(),
		RecreateDependencies: opts.dependenciesRecreateStrategy(),
		Inherit:              !opts.noInherit,
		Timeout:              opts.getTimeOut(),
		QuietPull:            opts.quietPull,
	}
	start := api.StartOptions{
		Wait:        opts.wait,
		WaitTimeout: opts.GetWaitTimeOut(),
	}
	if err := client.Up(ctx, create, start); err != nil {
		return errors.Wrap(err, "Execute docker compose up error")
	}

	// show notes
	var notesBuffer bytes.Buffer
	for fileName, content := range files {
		if strings.HasSuffix(fileName, app.DefaultNotesFileName) {
			// If buffer contains data, add newline before adding more
			if notesBuffer.Len() > 0 {
				notesBuffer.WriteString("\n")
			}
			notesBuffer.WriteString(content)
		}
	}
	if notesBuffer.Len() > 0 {
		fmt.Println(notesBuffer.String())
	}
	return nil
}
