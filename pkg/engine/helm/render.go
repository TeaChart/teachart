package helm

import (
	"path"
	_ "unsafe"

	"github.com/yp05327/teachart/pkg/values"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	helm_values "helm.sh/helm/v3/pkg/cli/values"
	helm_engine "helm.sh/helm/v3/pkg/engine"
)

type renderable struct {
	// tpl is the current template.
	tpl string
	// vals are the values to be supplied to the template.
	vals chartutil.Values
	// namespace prefix to the templates of the current chart
	basePath string
}

type files map[string][]byte

//go:linkname newFiles helm.sh/helm/v3/pkg/engine.newFiles
func newFiles(from []*chart.File) files

//go:linkname isTemplateValid helm.sh/helm/v3/pkg/engine.isTemplateValid
func isTemplateValid(ch *chart.Chart, templateName string) bool

//go:linkname render helm.sh/helm/v3/pkg/engine.(*Engine).render
func render(e *helm_engine.Engine, tpls map[string]renderable) (rendered map[string]string, err error)

// override helm.sh/helm/v3/pkg/engine.recAllTpls
func helm_engine_recAllTpls(c *chart.Chart, teachart *values.TeaChart, templates map[string]renderable, vals chartutil.Values) map[string]interface{} {
	subCharts := make(map[string]interface{})
	// override the chart meta data
	chartMetaData := struct {
		values.TeaChart
		IsRoot bool
	}{*teachart, c.IsRoot()}

	next := map[string]interface{}{
		"TeaChart":     chartMetaData,
		"Files":        newFiles(c.Files),
		"Release":      vals["Release"],
		"Capabilities": vals["Capabilities"],
		"Values":       make(chartutil.Values),
		"Subcharts":    subCharts,
	}

	// If there is a {{.Values.ThisChart}} in the parent metadata,
	// copy that into the {{.Values}} for this template.
	if c.IsRoot() {
		next["Values"] = vals["Values"]
	} else if vs, err := vals.Table("Values." + c.Name()); err == nil {
		next["Values"] = vs
	}

	for _, child := range c.Dependencies() {
		subCharts[child.Name()] = helm_engine_recAllTpls(child, teachart, templates, next)
	}

	newParentID := c.ChartFullPath()
	for _, t := range c.Templates {
		if t == nil {
			continue
		}
		if !isTemplateValid(c, t.Name) {
			continue
		}
		templates[path.Join(newParentID, t.Name)] = renderable{
			tpl:      string(t.Data),
			vals:     next,
			basePath: path.Join(newParentID, "templates"),
		}
	}

	return next
}

func (h *Helm) Render(valueOpts helm_values.Options, save bool) (map[string]string, error) {
	renderValues, err := h.getRenderValues(valueOpts)
	if err != nil {
		return nil, err
	}
	templates := make(map[string]renderable)
	helm_engine_recAllTpls(h.chart, h.teachart, templates, renderValues)
	files, err := render(h.engine, templates)
	if err != nil {
		return nil, err
	}
	if save {
		err = h.save(files)
	}
	return files, err
}
