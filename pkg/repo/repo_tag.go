// Copyright © 2024 TeaChart Authors

package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

// GetLatestTag returns latest tags in local
func (c *Client) GetLatestTag() (string, error) {
	if err := c.LoadGitRepo(); err != nil {
		return "", err
	}

	tags, err := c.gitRepo.Tags()
	if err != nil {
		return "", err
	}

	var latestCommit *object.Commit
	var latestTag string
	err = tags.ForEach(func(tag *plumbing.Reference) error {
		commit, err := c.getCommitFromTag(tag)
		if err != nil {
			return err
		}

		if latestCommit == nil || commit.Committer.When.After(latestCommit.Committer.When) {
			latestCommit = commit
			latestTag = tag.Name().String()
		}
		return nil
	})
	if err == nil && latestTag == "" {
		return "", git.ErrTagNotFound
	}
	return latestTag, err
}

// GetTag returns the tag in local
func (c *Client) GetTag(name string) (string, error) {
	if err := c.LoadGitRepo(); err != nil {
		return "", err
	}
	tag, err := c.gitRepo.Tag(name)
	if err != nil {
		return "", err
	}
	return tag.Name().String(), nil
}

// support for lightweight tag and annotated tag
func (c *Client) getCommitFromTag(tag *plumbing.Reference) (*object.Commit, error) {
	hash, err := c.gitRepo.ResolveRevision(plumbing.Revision(tag.Hash().String()))
	if err != nil {
		return nil, err
	}
	o, err := c.gitRepo.Object(plumbing.AnyObject, *hash)
	if err != nil {
		return nil, err
	}
	return o.(*object.Commit), nil
}

// GetTagByVersion returns the tag by version in local
func (c *Client) GetTagByVersion(version string) (tagRef string, err error) {
	if version == "" || version == "latest" {
		tagRef, err = c.GetLatestTag()
	} else {
		tagRef, err = c.GetTag(version)
	}
	logrus.Debugf("got tag ref `%s` from version `%s`", tagRef, version)
	return tagRef, err
}
