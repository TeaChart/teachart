// Copyright © 2024 TeaChart Authors

package repo

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	// git remote name
	remoteName string

	// repos folder path
	rootDir string
}

// NewManager returns a new manager
func NewManager(rootDir, remoteName string) *Manager {
	return &Manager{
		rootDir:    rootDir,
		remoteName: remoteName,
	}
}

// Init verify the root dir path and create the folder if not exist
func (m *Manager) Init() error {
	return os.MkdirAll(m.rootDir, 0755)
}

func (m *Manager) getRepoPath(name string) string {
	return filepath.Join(m.rootDir, name)
}

func (m *Manager) Remove(name string) error {
	err := os.RemoveAll(m.getRepoPath(name))
	if os.IsNotExist(err) {
		logrus.Warnf("repo `%s` does not exist in %s", name, m.rootDir)
		return nil
	}
	return err
}

func (m *Manager) List() (map[string]*git.Repository, error) {
	repos := make(map[string]*git.Repository)

	err := filepath.WalkDir(m.rootDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// find repos directory with depth 1
		if entry.IsDir() && strings.Count(strings.TrimPrefix(path, m.rootDir), string(os.PathSeparator)) == 1 {
			name := entry.Name()

			gitRepo, err := git.PlainOpen(path)
			if err != nil {
				logrus.Errorf("reading repo `%s`: %v", name, err)
				return nil
			}
			logrus.Debugf("found repo directory: %s", name)

			repos[name] = gitRepo
		}
		return nil
	})
	return repos, err
}

// Add download git repo to local
func (m *Manager) Add(ctx context.Context, name, url string, force bool) error {
	repoPath := m.getRepoPath(name)
	repoClient := NewClient(repoPath, NewClientOptions{
		RemoteName: m.remoteName,
	})

	// remove git repo if force
	if force {
		if err := os.RemoveAll(repoPath); err != nil {
			return err
		}
	}

	// load git repo
	fmt.Printf("adding repo `%s` from %s\n", name, url)
	if err := repoClient.Init(ctx, url); err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			// need rollback when init failed
			logrus.Errorf("init repo failed: %v", err)
			fmt.Print("try to rollback the changes: ")
			if err := os.RemoveAll(repoPath); err != nil {
				return err
			}
			fmt.Println("success.")
			return nil
		}
		return err
	}

	// TODO: verify teachart

	return nil
}

func (m *Manager) Checkout(ctx context.Context, name string, authMethod transport.AuthMethod, opts CheckOutOptions) error {
	repoPath := m.getRepoPath(name)
	repoClient := NewClient(repoPath, NewClientOptions{
		RemoteName: m.remoteName,
	})

	return repoClient.Checkout(ctx, authMethod, opts)
}
