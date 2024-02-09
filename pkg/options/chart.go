// Copyright © 2024 TeaChart Authors

package options

import (
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type ChartOptions struct {
	// Name is teachart name
	Name string
	// Version is the remote teachart version to use
	Version string
	// Commit is the remote teachart commit to use
	Commit string
	// URL is git url
	URL string
	// UserName is git login username
	UserName string
	// Password is git login password or privatekey passphrase
	Password string
	// PrivateKey is git login privatekey
	PrivateKey string
	// Dir is local dir to teachart, use for dev
	Dir string

	// rootDir is the folder saving all repos in local
	rootDir string
}

// NewChartOptions returns a new ChartOptions
func NewChartOptions(rootDir string) *ChartOptions {
	return &ChartOptions{
		rootDir: rootDir,
	}
}

// Get returns auth method
func (c *ChartOptions) GetAuthMethod() (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error

	if c.UserName != "" {
		auth = &http.BasicAuth{
			Username: c.UserName,
			Password: c.Password,
		}
	}

	if c.PrivateKey != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", c.PrivateKey, c.Password)
	}
	return auth, err
}

func (c *ChartOptions) GetChartDir() string {
	if c.Dir != "" {
		return c.Dir
	}
	return filepath.Join(c.rootDir, c.Name)
}
