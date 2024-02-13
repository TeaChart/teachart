// Copyright © 2024 TeaChart Authors

package engine

import (
	"helm.sh/helm/v3/pkg/cli/values"
)

type NewEngineOptions struct {
	Strict bool
}
type RenderEngine interface {
	Render(valueOpts values.Options, save bool) (map[string]string, error)
	GetConfigPaths(files map[string]string) []string
}

func NewRenderEngine(chartDir string, teachart *TeaChart, opts *NewEngineOptions) (RenderEngine, error) {
	return newHelm(chartDir, teachart, opts)
}
