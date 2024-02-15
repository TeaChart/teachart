package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

func TestInstallAndUninstall(t *testing.T) {
	ctx := context.Background()

	c, err := NewRootCmd(ctx, options.NewGlobalOptions(app.DefaultRepoDir))
	assert.NoError(t, err)

	tests := []string{
		"repo add test https://github.com/TeaChart/gitea",
		"install test",
		"uninstall",
		"repo remove test",
	}

	for _, tt := range tests {
		c.SetArgs(strings.Split(tt, " "))
		assert.NoError(t, c.Execute())
	}
}
