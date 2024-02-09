// Copyright © 2024 TeaChart Authors

package options

type RepoListOptions struct {
	*RepoOptions
}

func NewRepoListOptions(repoOptions *RepoOptions) *RepoListOptions {
	return &RepoListOptions{
		RepoOptions: repoOptions,
	}
}
