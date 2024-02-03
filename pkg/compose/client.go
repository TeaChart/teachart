// Copyright © 2024 TeaChart Authors

package compose

import (
	"context"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	docker_cmd "github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

type Client struct {
	*compose_cmd.ProjectOptions

	dockerCli      *docker_cmd.DockerCli
	composeService api.Service
	services       []string
}

// NewClient returns a new docker compose client
// projectOpts.ConfigPaths should be provided
func NewClient(services []string, projectOpts *compose_cmd.ProjectOptions) (*Client, error) {
	dockerCli, err := docker_cmd.NewDockerCli()
	if err != nil {
		return nil, err
	}
	if err := dockerCli.Initialize(&flags.ClientOptions{
		Context: dockerCli.CurrentContext(),
	}); err != nil {
		return nil, err
	}

	return &Client{
		ProjectOptions: projectOpts,
		dockerCli:      dockerCli,
		composeService: compose.NewComposeService(dockerCli),
		services:       services,
	}, nil
}

func (c *Client) toProject(pofs []cli.ProjectOptionsFn) (*types.Project, error) {
	return c.ProjectOptions.ToProject(c.dockerCli, c.services, pofs...)
}

func (c *Client) Up(ctx context.Context, createOpts api.CreateOptions, startOpts api.StartOptions) error {
	project, err := c.toProject([]cli.ProjectOptionsFn{
		cli.WithResolvedPaths(true),
		cli.WithDiscardEnvFile,
		cli.WithContext(ctx),
	})
	if err != nil {
		return err
	}

	startOpts.Project = project
	startOpts.Services = c.services
	createOpts.Services = c.services

	return c.composeService.Up(ctx, project, api.UpOptions{
		Create: createOpts,
		Start:  startOpts,
	})
}

func (c *Client) Down(ctx context.Context, downOpts api.DownOptions) error {
	downOpts.Services = c.services
	return c.composeService.Down(ctx, c.ProjectName, downOpts)
}

func (c *Client) Config(ctx context.Context, pofs []cli.ProjectOptionsFn, configOpts api.ConfigOptions) ([]byte, error) {
	project, err := c.toProject([]cli.ProjectOptionsFn{
		cli.WithResolvedPaths(true),
		cli.WithDiscardEnvFile,
		cli.WithContext(ctx),
	})
	if err != nil {
		return nil, err
	}
	return c.composeService.Config(ctx, project, configOpts)
}
