// Copyright © 2024 TeaChart Authors

package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/yp05327/teachart/pkg/app"
)

type Client struct {
	// necessary params
	path       string
	remoteName string

	gitRepo *git.Repository
}

type NewClientOptions struct {
	// git remote name, can be empty
	RemoteName string
}

// NewClient create a new repo client
// path should be the full path to the repo
func NewClient(path string, opts NewClientOptions) *Client {
	if opts.RemoteName == "" {
		opts.RemoteName = app.DefaultRemoteName
	}

	return &Client{
		path:       path,
		remoteName: opts.RemoteName,
	}
}

// Path returns the full path of the repo
func (c *Client) Path() string {
	return c.path
}

// Name returns the name of the repo
func (c *Client) Name() string {
	return filepath.Base(c.path)
}

// Remove removes the git repo saved in local
func (c *Client) Remove() error {
	return os.RemoveAll(c.path)
}

// LoadGitRepoWithInit loads the git repo from local
// if opts is not nil, init the directory when not exist
func (c *Client) LoadGitRepo() (err error) {
	if c.gitRepo != nil {
		return nil
	}
	if c.path == "" {
		return fmt.Errorf("load repo failed: repo path is empty")
	}
	c.gitRepo, err = git.PlainOpenWithOptions(c.path, &git.PlainOpenOptions{})
	return err
}

// GetRemote returns the remote of the repo
func (r *Client) GetRemote() (*git.Remote, error) {
	if err := r.LoadGitRepo(); err != nil {
		return nil, err
	}
	return r.gitRepo.Remote(r.remoteName)
}
