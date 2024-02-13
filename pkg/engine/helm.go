// Copyright © 2024 TeaChart Authors

package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yp05327/teachart/pkg/app"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	helm_engine "helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/getter"
)

type Helm struct {
	engine *helm_engine.Engine
	chart  *chart.Chart

	teachart *TeaChart
}

func newHelm(chartDir string, teachart *TeaChart, opts *NewEngineOptions) (*Helm, error) {
	chart, err := loader.Load(chartDir)
	if err != nil {
		return nil, err
	}
	// TODO support Dependencies

	if opts == nil {
		opts = &NewEngineOptions{}
	}
	return &Helm{
		engine: &helm_engine.Engine{
			Strict: opts.Strict,
		},
		chart:    chart,
		teachart: teachart,
	}, nil
}

func (h *Helm) Render(valueOpts values.Options, save bool) (map[string]string, error) {
	renderValues, err := h.getRenderValues(valueOpts)
	if err != nil {
		return nil, err
	}
	files, err := h.engine.Render(h.chart, renderValues)
	if err != nil {
		return nil, err
	}
	if save {
		err = h.save(files)
	}
	return files, err
}

func (h *Helm) getRenderValues(valueOpts values.Options) (chartutil.Values, error) {
	// user define values
	userValues, err := valueOpts.MergeValues(getter.Providers{})
	if err != nil {
		return nil, err
	}

	top := map[string]interface{}{
		"Chart": h.chart.Metadata,
	}

	userValues["TeaChart"] = h.teachart
	values, err := chartutil.CoalesceValues(h.chart, userValues)
	if err != nil {
		return top, err
	}

	if err := chartutil.ValidateAgainstSchema(h.chart, values); err != nil {
		errFmt := "values don't meet the specifications of the schema(s) in the following chart(s):\n%s"
		return top, fmt.Errorf(errFmt, err.Error())
	}
	top["Values"] = values
	return top, nil

}

func (h *Helm) save(files map[string]string) error {
	for fileName, content := range files {
		if len(content) == 0 {
			continue
		}

		path := h.getSavePath(fileName)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		perm := fs.FileMode(0644)
		if filepath.Ext(path) == ".sh" {
			perm = 0755
		}
		if err := os.WriteFile(path, []byte(content), perm); err != nil {
			return err
		}
	}
	return nil
}

func (h *Helm) GetConfigPaths(files map[string]string) []string {
	paths := make([]string, 0, len(files))
	for fileName, data := range files {
		if len(data) == 0 {
			continue
		}
		path := h.getSavePath(fileName)
		ext := filepath.Ext(path)
		if ext == ".yaml" || ext == ".yml" {
			paths = append(paths, path)
		}
	}
	return paths
}

func (h *Helm) getSavePath(fileName string) string {
	absPath := strings.Replace(fileName, path.Join(h.chart.Name(), "templates"), "", 1)
	return filepath.Join(h.teachart.ProjectDir, app.DefaultTempDirName, absPath)
}
