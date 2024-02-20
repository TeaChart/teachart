// Copyright © 2024 TeaChart Authors

package engine

import (
	"github.com/yp05327/teachart/pkg/engine/helm"
	"github.com/yp05327/teachart/pkg/values"
	helm_values "helm.sh/helm/v3/pkg/cli/values"
)

type RenderEngine interface {
	Render(valueOpts helm_values.Options, save bool) (map[string]string, error)
	GetConfigPaths(files map[string]string) []string
}

func NewRenderEngine(chartDir string, teachart *values.TeaChart, strict bool) (RenderEngine, error) {
	return helm.New(chartDir, teachart, strict)
}
