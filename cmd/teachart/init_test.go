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

func TestInit(t *testing.T) {
	ctx := context.Background()

	var out bytes.Buffer
	c := NewInitCmd(ctx, options.NewGlobalOptions(app.DefaultRepoDir))
	c.SetOut(&out)
	c.SetArgs(strings.Split("init --name test", " "))
	assert.NoError(t, c.Execute())
	assert.Equal(t, "test created\n", out.String())
}
