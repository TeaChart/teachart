// Copyright © 2024 TeaChart Authors

package options

type RepoRemoveOptions struct {
	*RepoOptions
	Name string
}

func NewRepoRemoveOptions(repoOptions *RepoOptions) *RepoRemoveOptions {
	return &RepoRemoveOptions{
		RepoOptions: repoOptions,
	}
}
