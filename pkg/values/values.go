// Copyright © 2024 TeaChart Authors

package values

import "helm.sh/helm/v3/pkg/chart"

// TeaChart is an extended chart metadata
type TeaChart struct {
	*chart.Metadata

	ProjectName string
	ProjectDir  string
	TempDir     string
}
