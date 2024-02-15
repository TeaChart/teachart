// Copyright © 2024 TeaChart Authors

package app

const (
	// teachart
	DefaultRepoDir      string = "repos"
	DefaultTemplatesDir string = "templates"

	DefaultRemoteName string = "origin"

	DefaultMetadataFileName     string = "Chart.yaml"
	DefaultValuesFileName       string = "values.yaml"
	DefaultValuesSchemaFileName string = "values.schema.json"
	DefaultNotesFileName        string = "NOTES.txt"

	// docker compose
	DefaultComposeExecutable string = "docker"
	DefaultTempDirName       string = ".teachart"
)

const DefaultMetadataFile string = `name: %s
description: A TeaChart simple
version: 0.0.1
appVersion: "0.0.1"
`

const DefaultValuesFile string = `# Default values for %s.
# This is a yaml file
# You can use these values in your templates by {{ .Values.XXX }}.
`

const DefaultNotesFile string = `# Notes for %s.
# This will display when the installation is finished.
# This file will be rendered by go template.
`
