YAML Include Transformer
========================

A simple YAML processor that implements include directives for YAML files.

## Example of Standalone Usage

An example of YAML input:


```yaml
program:
  language: lua
  code!textfile:  source.lua
  data!base64file: data.bin
  labels:jsonfile: labels.json

```

Run yaml processor:

```shell
yaml-include-transformer < test.yaml
```

Output:

```yaml
program:
  code: |
    print("Hello!\n")
  data: ODIzY2YxODYyNDVmNTBkMzk0YjMxMDlmYTNiM2E5NjYgIC0K
  labels:
	foo: bar
	results: [1, 2, 3]
  language: lua
```

## Usage as kustomize plugin

Installation as an "exec" plugin:

```shell
yaml-include-transformer -i
Installing kustomize exec plugin /home/dudin/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer
copy './yaml-include-transformer' to '/home/dudin/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer/FieldIncludeTransformer'
Installation complete
```

Create plugin configuration file (p.e. include-plugin.yaml):

```shell
yaml-include-transformer -p > include-plugin.yaml

```

Add a transformer declaration to the `kustomization.yaml` file:

```yaml
transformers:
  - include-plugin.yaml
```
Invoke kustomize build:

```shell
kustomize --enable-exec --enable-alpha-plugins build 
```

## Command Line Arguments Reference


Usage: 

```
  yaml-include-transformer [configfile] [options ...]
```
Options:

* `-h --help`	        Print usage message
* `-i --install`        Install as kustomize exec plugin
* `-p --plugin-conf`    Print kustomize plugin configuration file
* `-f --file file.yaml` Specify input file instead of standard input
* `-u --updir`          Allow specifying .. in file paths
* `-l --links`          Allow following symlinks in file paths
* `-a --abs`            Allow absolute paths in file paths
* `-v --version`        Print program version


## Supported Include directives

* `foo!textfile: file.txt`    include file.txt as a text field
* `bar!base64file: file.bin`  include file.bin as base64 text


## Configuration File

Accepting configuration file as first program argument is required for
compatibility with the Kustomize exec plugin protocol. The
configuration file is accepted but not actually used.

## Usage as Kustomize shared library based plugin

[TBD]

## Usage as Kustomize KRM function based plugin

[TBD]

## Configuring ArgoCD to use the plugin

[TBD]





