package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yp05327/teachart/pkg/app"
	"helm.sh/helm/v3/pkg/cli/values"
)

const testChartDir = "../../tests/data/chart"
const testValuesFile = "../../tests/data/values.yaml"

var testTeaChart = TeaChart{
	ProjectName: "test",
	ProjectDir:  "",
	TempDir:     app.DefaultTemplatesDir,
}

func TestHelmRender(t *testing.T) {
	e, err := NewRenderEngine(testChartDir, &testTeaChart, &NewEngineOptions{Strict: true})
	assert.NoError(t, err)

	_, err = e.Render(values.Options{
		ValueFiles: []string{testValuesFile},
	}, false)
	assert.NoError(t, err)
}
