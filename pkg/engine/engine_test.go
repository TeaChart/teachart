package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/values"
	helm_values "helm.sh/helm/v3/pkg/cli/values"
)

const testChartDir = "../../tests/data/chart"
const testValuesFile = "../../tests/data/values.yaml"

var testTeaChart = values.TeaChart{
	ProjectName: "test",
	ProjectDir:  "",
	TempDir:     app.DefaultTemplatesDir,
}

func TestHelmRender(t *testing.T) {
	e, err := NewRenderEngine(testChartDir, &testTeaChart, true)
	assert.NoError(t, err)

	_, err = e.Render(helm_values.Options{
		ValueFiles: []string{testValuesFile},
	}, false)
	assert.NoError(t, err)
}
