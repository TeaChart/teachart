// Copyright © 2024 TeaChart Authors

package options

type RepoAddOptions struct {
	*RepoOptions
	*ChartOptions

	Name  string
	Force bool
}

func NewRepoAddOptions(repoOptions *RepoOptions) *RepoAddOptions {
	return &RepoAddOptions{
		RepoOptions:  repoOptions,
		ChartOptions: NewChartOptions(repoOptions.GetRepoRootDir()),
	}
}
