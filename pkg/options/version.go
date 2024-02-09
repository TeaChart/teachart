// Copyright © 2024 TeaChart Authors

package options

type VersionOptions struct {
	*GlobalOptions

	Short  bool
	Format string
}

func NewVersionOptions(globalOptions *GlobalOptions) *VersionOptions {
	return &VersionOptions{
		GlobalOptions: globalOptions,
	}
}
