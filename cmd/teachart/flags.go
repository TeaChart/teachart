// Copyright © 2024 TeaChart Authors

package cmd

import (
	"os"

	"github.com/compose-spec/compose-go/v2/utils"
	compose_cmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/yp05327/teachart/pkg/options"
	"helm.sh/helm/v3/pkg/cli/values"
)

func addValueFlags(f *pflag.FlagSet, v *values.Options) {
	f.StringSliceVarP(&v.ValueFiles, "values", "f", []string{}, "values file in yaml, can specify multiple files")
	f.StringArrayVar(&v.Values, "set", []string{}, "set values on the command line, can specify multiple or separate values with commas: key1=val1,key2=val2")
	f.StringArrayVar(&v.StringValues, "set-string", []string{}, "set string values on the command line, can specify multiple or separate values with commas: key1=string1,key2=string2")
	f.StringArrayVar(&v.FileValues, "set-file", []string{}, "set values from respective files specified via the command line, can specify multiple or separate values with commas: key1=path1,key2=path2")
	f.StringArrayVar(&v.JSONValues, "set-json", []string{}, "set json values on the command line, can specify multiple or separate values with commas: key1=json1,key2=json2")
	f.StringArrayVar(&v.LiteralValues, "set-literal", []string{}, "set a literal string value on the command line")

}

func addChartFlags(f *pflag.FlagSet, c *options.ChartOptions) {
	f.StringVarP(&c.Version, "version", "v", "latest", "specify a version for the chart version to use. default is `latest`")
	f.StringVar(&c.Commit, "commit", "", "specify a commit to use. if this flag is set, --version/-v will be ignored")
	f.StringVar(&c.URL, "repo", "", "chart repository url where to locate the requested chart")
	f.StringVarP(&c.UserName, "username", "u", "", "username for chart repo")
	f.StringVar(&c.Password, "password", "", "password for chart repo or the privatekey passphrase")
	f.StringVar(&c.PrivateKey, "private-key", "", "private key for chart repo")
	f.StringVar(&c.Dir, "dir", "", "specify a dir to teachart. will ignore the given chart name")
}

func addCreateFlags(f *pflag.FlagSet, c *createOptions) {
	f.StringVar(&c.Pull, "pull", "policy", `Pull image before running ("always"|"missing"|"never")`)
	f.BoolVar(&c.removeOrphans, "remove-orphans", false, "Remove containers for services not defined in the Compose file.")
	f.StringArrayVar(&c.scale, "scale", []string{}, "Scale SERVICE to NUM instances. Overrides the `scale` setting in the Compose file if present.")
	f.BoolVar(&c.forceRecreate, "force-recreate", false, "Recreate containers even if their configuration and image haven't changed.")
	f.BoolVar(&c.noRecreate, "no-recreate", false, "If containers already exist, don't recreate them. Incompatible with --force-recreate.")
	f.IntVarP(&c.timeout, "timeout", "t", 0, "Use this timeout in seconds for container shutdown when attached or when containers are already running.")
	f.BoolVar(&c.recreateDeps, "always-recreate-deps", false, "Recreate dependent containers. Incompatible with --no-recreate.")
	f.BoolVarP(&c.noInherit, "renew-anon-volumes", "V", false, "Recreate anonymous volumes instead of retrieving data from the previous containers.")
	f.BoolVar(&c.quietPull, "quiet-pull", false, "Pull without printing progress information.")
}

func addStartFlags(f *pflag.FlagSet, s *startOptions) {
	f.BoolVar(&s.noColor, "no-color", false, "Produce monochrome output.")
	f.BoolVar(&s.noPrefix, "no-log-prefix", false, "Don't print prefix in logs.")
	f.BoolVar(&s.noStart, "no-start", false, "Don't start the services after creating them.")
	f.BoolVar(&s.timestamp, "timestamps", false, "Show timestamps.")
	f.BoolVar(&s.noDeps, "no-deps", false, "Don't start linked services.")
	f.BoolVar(&s.wait, "wait", false, "Wait for services to be running|healthy. Implies detached mode.")
	f.IntVar(&s.waitTimeout, "wait-timeout", 0, "Maximum duration to wait for the project to be running|healthy.")
}

func addDownFlags(f *pflag.FlagSet, d *downOptions) {
	removeOrphans := utils.StringToBool(os.Getenv(compose_cmd.ComposeRemoveOrphans))
	f.BoolVar(&d.removeOrphans, "remove-orphans", removeOrphans, "Remove containers for services not defined in the Compose file.")
	f.IntVarP(&d.timeout, "timeout", "t", 0, "Specify a shutdown timeout in seconds")
	f.BoolVarP(&d.volumes, "volumes", "v", false, `Remove named volumes declared in the "volumes" section of the Compose file and anonymous volumes attached to containers.`)
	f.StringVar(&d.images, "rmi", "", `Remove images used by services. "local" remove only images that don't have a custom tag ("local"|"all")`)
	f.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		if name == "volume" {
			name = "volumes"
			logrus.Warn("--volume is deprecated, please use --volumes")
		}
		return pflag.NormalizedName(name)
	})
}

func addConfigFlasg(f *pflag.FlagSet, c *configOptions) {
	f.StringVar(&c.Format, "format", "yaml", "Format the output. Values: [yaml | json]")
	f.BoolVar(&c.resolveImageDigests, "resolve-image-digests", false, "Pin image tags to digests.")
	f.BoolVar(&c.noInterpolate, "no-interpolate", false, "Don't interpolate environment variables.")
	f.BoolVar(&c.noNormalize, "no-normalize", false, "Don't normalize compose model.")
	f.BoolVar(&c.noResolvePath, "no-path-resolution", false, "Don't resolve file paths.")
	f.BoolVar(&c.noConsistency, "no-consistency", false, "Don't check model consistency - warning: may produce invalid Compose output")
}
