package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

func TestLint(t *testing.T) {
	ctx := context.Background()

	var out bytes.Buffer
	c := NewLintCmd(ctx, options.NewGlobalOptions(app.DefaultRepoDir))
	c.SetOut(&out)
	c.SetArgs(strings.Split("../../tests/data/chart -f ../../tests/data/values.yaml", " "))
	assert.NoError(t, c.Execute())
	assert.Contains(t, out.String(), "1 chart(s) linted, 0 chart(s) failed")
}
