// Copyright © 2024 TeaChart Authors

package repo

import (
	"context"
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type FetchOptions struct {
	Ref   config.RefSpec
	Debug bool
}

// Fetch fetch the repo from remote, you can assign a special tag in options
func (c *Client) Fetch(ctx context.Context, authMethod transport.AuthMethod, opts FetchOptions) error {
	if err := c.LoadGitRepo(); err != nil {
		return err
	}

	gitOpts := &git.FetchOptions{
		RemoteName: c.remoteName,
		Auth:       authMethod,
		Depth:      1,
	}
	if opts.Ref != "" {
		gitOpts.RefSpecs = []config.RefSpec{opts.Ref}
	}
	if opts.Debug {
		gitOpts.Progress = os.Stdout
	}
	err := c.gitRepo.FetchContext(ctx, gitOpts)
	// ignore already up-to-date err
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}
