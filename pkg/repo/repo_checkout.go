// Copyright © 2024 TeaChart Authors

package repo

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type CheckOutOptions struct {
	Version string
	Commit  string
	Debug   bool
}

// Checkout fetch all tags from remote and checout the git repo by given version in the options
func (c *Client) Checkout(ctx context.Context, authMethod transport.AuthMethod, opts CheckOutOptions) error {
	gitOpts := &git.CheckoutOptions{
		Force: true,
	}
	fetchOpts := FetchOptions{
		Debug: opts.Debug,
	}

	if opts.Commit != "" {
		gitOpts.Hash = plumbing.NewHash(opts.Commit)
	} else {
		tag, err := c.GetTagByVersion(opts.Version)
		if err != nil {
			return err
		}
		fetchOpts.Ref = config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag))
		gitOpts.Branch = plumbing.ReferenceName(tag)
	}

	if err := c.Fetch(ctx, authMethod, fetchOpts); err != nil {
		return err
	}

	worktree, err := c.gitRepo.Worktree()
	if err != nil {
		return err
	}
	return worktree.Checkout(gitOpts)
}
