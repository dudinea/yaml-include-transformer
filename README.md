YAML Include Transformer
========================

A simple YAML processor that implements include directives for YAML
files. It can be used as a standalone utility as well as a plugin for
[Kustomize](https://kustomize.io).

## Standalone Usage

An example of YAML input:

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: demo-cm
  labels!jsonfile: labels.json
  annotations!yamlfile: annotations.yaml
data:
  language: lua
  code!textfile:  source.lua
  data!base64file: data.bin
```

Run yaml processor:

```shell
yaml-include-transformer < examples.yaml
```

Output:

```yaml
---
apiVersion: v1
data:
    code: |
        print("Hello!\n")
    data: hczjkOrano3o4Womxt0SFtxXVo4MuSph4w==
    language: lua
kind: ConfigMap
metadata:
    annotations:
        aprefix/akey: avalue
    labels:
        app: demo
        environment: dev
    name: demo-cm
```

## Usage as Kustomize Plugin

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

* `foo!textfile: file.txt`    include `file.txt` as a text field.
* `foo!base64file: file.bin`  include `file.bin` as base64 text.
* `foo!jsonfile: file.json`   deserialize `file.json` and include it as data structure.
* `foo!yamlfile: file.yaml`   deserialize `file.yaml` and include it as data structure.


## Configuration File

Accepting configuration file as first program argument is required for
compatibility with the Kustomize exec plugin protocol. The
configuration file is accepted but currently it is not actually used.

## Configuring ArgoCD to use Kustomize with the plugin 

[ArgoCD](https://argoproj.github.io) can be customized to support
YAML Include Transformer with kustomize-based applications. One needs
to modify `argocd-repo-server` deployment to use a customized
docker image and to change kustomize command line flags.  `kustomize.buildOptions` 
in the `argocd-cm` ConfigMap. 

### Building a Customized ArgoCD Image

This command will add the `yaml-include-transformer` binary to the
source argocd docker image and installs it as a customize plugin.  You
can customize target repository and source image using environment
variables, see details in the Makefile.


```shell
$ env ARGOCD_REPO=some-repo/argocd-yit ARGOCD_VER=v2.4.4  make argo_docker_build
echo 	"FROM quay.io/argoproj/argocd:v2.4.4 \n" \
	"ADD ./yaml-include-transformer /usr/local/bin\n" \
	"RUN /usr/local/bin/yaml-include-transformer -i\n" > Dockerfile.argocd
docker build -f Dockerfile.argocd -t some-repo/argocd-yit:v2.4.4_yitv0.0.2-alpha .
Sending build context to Docker daemon  7.269MB
Step 1/3 : FROM quay.io/argoproj/argocd:v2.4.4
 ---> 34842ba61a5a
Step 2/3 : ADD ./yaml-include-transformer /usr/local/bin
 ---> Using cache
 ---> 4a2f7c58907e
Step 3/3 : RUN /usr/local/bin/yaml-include-transformer -i
 ---> Using cache
 ---> f3c91076e12e
Successfully built f3c91076e12e
Successfully tagged some-repo/argocd-yit:v2.4.4_yitv0.0.2-alpha
```

`make argo_docker_push` will push the image to the repository.

### Patching the ArgoCD configuration

The following command patches the deployment of `argocd-repo-server` to use the customized
docker image and changes the kustomize command line flags in the parameter `kustomize.buildOptions` 
in the `argocd-cm` ConfigMap. See details in the Makefile.

```shell
 $ /usr/bin/env ARGOCD_REPO=some-repo/argocd-yit ARGOCD_VER=v2.4.4  make argo_patch
kubectl patch deployment -n  argocd argocd-repo-server -p \
'{"spec" : {"template" : { "spec" : { "containers" : [ { "image" : "some-repo/argocd-yit:v2.4.4_yitv0.0.2-alpha", "name" : "argocd-repo-server"  }]}}}}'
deployment.apps/argocd-repo-server patched
kubectl patch cm -n argocd argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-exec --enable-alpha-plugins"}}'
configmap/argocd-cm patched
```

## Usage as Kustomize shared library based plugin

[TBD]

## Usage as Kustomize KRM function based plugin

[TBD]







