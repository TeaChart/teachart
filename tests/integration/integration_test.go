package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	cmd "github.com/yp05327/teachart/cmd/teachart"
	"github.com/yp05327/teachart/pkg/app"
	"github.com/yp05327/teachart/pkg/options"
)

func TestRepo(t *testing.T) {
	ctx := context.Background()

	var out bytes.Buffer
	c := cmd.NewRepoCmd(ctx, options.NewGlobalOptions(app.DefaultRepoDir))
	c.SetOut(&out)

	tests := []struct {
		args     string
		expected string
	}{
		{
			// list without repos exist
			args:     "list",
			expected: "no repositories. use `repo add` to add repos",
		},
		{
			// add a repo
			args:     "add test https://github.com/TeaChart/gitea",
			expected: "adding repo `test` from https://github.com/TeaChart/gitea",
		},
		{
			// list repo
			args:     "list",
			expected: "NAME\tURL                              \ntest\thttps://github.com/TeaChart/gitea",
		}, {
			// remove the repo
			args:     "remove test",
			expected: "repo test removed",
		},
	}

	for _, tt := range tests {
		c.SetArgs(strings.Split(tt.args, " "))
		err := c.Execute()
		assert.NoError(t, err)
		length := len(out.String())
		assert.Less(t, 0, length)
		assert.Equal(t, tt.expected, out.String()[:length-1])
		out.Reset()
	}
}

func TestInstallAndUninstall(t *testing.T) {
	ctx := context.Background()

	c, err := cmd.NewRootCmd(ctx, options.NewGlobalOptions(app.DefaultRepoDir))
	assert.NoError(t, err)

	tests := []string{
		"repo add test https://github.com/TeaChart/gitea",
		"install test",
		"uninstall",
	}

	for _, tt := range tests {
		c.SetArgs(strings.Split(tt, " "))
		err := c.Execute()
		assert.NoError(t, err)
	}
}
