// Copyright © 2024 TeaChart Authors

package options

type InitOptions struct {
	*GlobalOptions

	Name string
}

func NewInitOptions(globalOptions *GlobalOptions) *InitOptions {
	return &InitOptions{
		GlobalOptions: globalOptions,
	}
}
