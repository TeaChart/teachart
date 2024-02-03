// Copyright © 2024 TeaChart Authors

package options

import "github.com/yp05327/teachart/pkg/repo"

type RepoOptions struct {
	*GlobalOptions
	Manager *repo.Manager
}

func NewRepoOptions(globalOptions *GlobalOptions) *RepoOptions {
	return &RepoOptions{
		GlobalOptions: globalOptions,
	}
}
