// Copyright © 2024 TeaChart Authors

package options

import (
	"os"
	"path/filepath"

	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/yp05327/teachart/pkg/values"
)

// GlobalOptions is the global configuration
type GlobalOptions struct {
	// Executable is the path to the docker compose
	ComposeExecutable string

	// ProjectDir is the project dir to the workspace
	ProjectDir string
	// ProjectName is the project name of docker compose
	ProjectName string
	// TempDirName is the folder name of saving generated docker compose files
	TempDirName string
	// Quiet is true if the output should be quiet
	Quiet bool
	// Debug is true if the output should be verbose
	Debug bool
	// LogLevel is the log level to use.
	LogLevel string

	// reposDirName is repos folder name.
	// default is app.DefaultRepoDir and can not be changed now.
	reposDirName string
}

// NewGlobalOptions return a new globaloptions
func NewGlobalOptions(reposDirName string) *GlobalOptions {
	return &GlobalOptions{
		reposDirName: reposDirName,
	}
}

func (g *GlobalOptions) GetInstallDir() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(path)
}

// GetExecutable returns the no color flag
func (g *GlobalOptions) GetComposeExecutable() string {
	return g.ComposeExecutable
}

// GetProjectDir returns the projectdir flag
func (g *GlobalOptions) GetProjectDir() string {
	projectDir := g.ProjectDir
	if projectDir == "" {
		projectDir, _ = os.Getwd()
	}
	return projectDir
}

// GetProjectName returns the project name flag
func (g *GlobalOptions) GetProjectName() string {
	if g.ProjectName != "" {
		return g.ProjectName
	}
	return filepath.Base(g.GetProjectDir())
}

// GetTempDir returns the temp dir path
func (g *GlobalOptions) GetTempDir() string {
	return filepath.Join(g.GetProjectDir(), g.TempDirName)
}

// GetProjectName returns the project name flag
func (g *GlobalOptions) GetRepoRootDir() string {
	return filepath.Join(g.GetInstallDir(), g.reposDirName)
}

func (g *GlobalOptions) GetProjectOptions() *compose_cmd.ProjectOptions {
	return &compose_cmd.ProjectOptions{
		ProjectName: g.GetProjectName(),
		ProjectDir:  g.GetProjectDir(),
	}
}

func (g *GlobalOptions) GetTeaChart() *values.TeaChart {
	return &values.TeaChart{
		ProjectName: g.GetProjectName(),
		ProjectDir:  g.GetProjectDir(),
		TempDir:     g.GetTempDir(),
	}
}
