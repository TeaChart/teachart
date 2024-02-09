// Copyright © 2024 TeaChart Authors

package options

import "strings"

type ComposeOptions struct {
	// Args is the list of arguments to pass to the the docker compose
	Args string
}

// NewComposeOptions returns a new ComposeOptions
func NewComposeOptions() *ComposeOptions {
	return &ComposeOptions{}
}

// GetArgs returns the args flag
func (c *ComposeOptions) GetArgs() []string {
	return strings.Split(c.Args, " ")
}
