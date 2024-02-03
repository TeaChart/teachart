# TeaChart

TeaChart is a tool for managing docker containers like Helm Chart.
As usual, we use [docker compose](https://github.com/docker/compose) to manage our docker containers.
But for some big systems which have lots of containers, we devote our energy in maintaining config files.
Then TeaChart can help to do this.
TeaChart uses [Helm Chart](https://github.com/helm/helm)'s render engine to generate config files automatically.
This project is WIP. May have breaking changes in the future.

## Install

As project is WIP, so binary releases will not be provided now.
We recommend you using [Dev Container](https://code.visualstudio.com/docs/devcontainers/containers), then there's no need to prepare the build enviroment.

Then run the command below in terminal to build the binary.

``` bash
make build
```

At last, you can add the binary `teachart` to your `PATH`, then you are ready to go.

## Usage

### For user

First, you need to add a repository. All repositories are managed by `go-git`.
They will be all saved in `repos` folder under the binary install folder.

``` bash
# add teachart repository, NAME will be the folder name. 
# So disallowed charactors in folder name, can not be used here
teachart repo add REPO_NAME https://github.com/xxx/xxx
# if the repo folder already exists, you can add --force/-f flag to overwirte it.
teachart repo add REPO_NAME https://github.com/xxx/xxx -f
```

Then you can install the added chart repo.
But before do that, you should notice that there are many changeable values in a chart.
All of them have default values. For most of the time, the default values should be overwirted.
For example, the database username and password. So if necessary, you can edit your own `values.yaml`.
And pass it by using `--values/-f` flag as following.

``` bash
teachart install REPO_NAME -f PATH_TO_VALUES_YAML
```

Sometimes, chart may not provide changeable values.
You can use `--set` flag to change the docker compose config file, what ever you want.

``` bash
teachart install REPO_NAME --set xxx.xxx.xxx=xxx
```

After the installation, teachart will generate docker compose yaml files in `.teachart` under the work directory.
Ensure not deleting `.teachart` folder, as it will be used when uninstall the chart.
You can uninstall the chart by the following command.

``` bash
teachart uninstall
```

You can use `-h` flag to check more optional flags which are not mentioned here.

### For Chart developer

As we use Helm's engine, you can use all features from Helm.

And we added `.Values.TeaChart` for some special values from TeaChart.

### For TeaChart developer

TODO
