package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yp05327/teachart/pkg/app"
	"helm.sh/helm/v3/pkg/cli/values"
)

const testChartDir = "../../tests/chart"

var testTeaChart = TeaChart{
	ProjectName: "test",
	ProjectDir:  "",
	TempDir:     app.DefaultTemplatesDir,
}

func TestHelmRender(t *testing.T) {
	e, err := NewRenderEngine(testChartDir, &testTeaChart)
	assert.NoError(t, err)

	_, err = e.Render(values.Options{}, false)
	assert.NoError(t, err)
}
