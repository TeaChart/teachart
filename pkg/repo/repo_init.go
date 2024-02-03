// Copyright © 2024 TeaChart Authors

package repo

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Init inits the git repo, and fetch all tags from remote
func (c *Client) Init(ctx context.Context, url string) error {
	if c.path == "" {
		return fmt.Errorf("init repo failed: repo path is empty")
	}

	var err error
	c.gitRepo, err = git.PlainInit(c.path, false)
	if err != nil {
		return err
	}

	remote, err := c.gitRepo.CreateRemote(&config.RemoteConfig{
		Name:   c.remoteName,
		URLs:   []string{url},
		Mirror: false,
	})
	if err != nil {
		return err
	}
	return remote.Config().Validate()
}
