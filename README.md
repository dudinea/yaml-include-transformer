YAML Include Transformer
========================

A simple YAML processor that implements include directives for YAML
files. It can be used as a standalone utility as well as a plugin for
[Kustomize](https://kustomize.io) and with [ArgoCD](https://argoproj.github.io).

Standalone Usage
----------------

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

Command Line Arguments Reference
--------------------------------

Usage: 

```
yaml-include-transformer [configfile] | [options ...]
```
Options:
* `-h --help`	           Print this usage message
* `-i --install`           Install as kustomize exec plugin
* `-p --plugin-conf`       Print kustomize plugin configuration file
* `-E --exec`              Exec plugin (for -p and -i)
* `-L --legacy`            Legacy  plugin (for -p and -i), default
* `-K --krm`               KRM-function plugin (for -p and -i)
* `-D --dockertag`         KRM-function docker tag
* `-f --file file.yaml ..` Specify Input files
* `-u --up-dir`            Allow specifying .. in file paths
* `-l --links`             Allow following symlinks in file paths
* `-a --abs`               Allow absolute paths in file paths
* `-v --version`           Print program version
* `-d --debug`             Print debug messages on stderr

Supported Include directives
----------------------------

* `foo!textfile: file.txt`    include `file.txt` as a text field.
* `foo!base64file: file.bin`  include `file.bin` as base64 text.
* `foo!jsonfile: file.json`   deserialize `file.json` and include it as data structure.
* `foo!yamlfile: file.yaml`   deserialize `file.yaml` and include it as data structure.

Usage as Kustomize Plugin
-------------------------

[Kustomize](https://kustomize.io) offers a plugin framework that
allows to add user defined transformers that make changes to existing
Kubernetes objects. Generally speaking, transformers get YAML
multi-document as their standard input, transform it in some way, and print it on their standard output.

The Kustomize plugins are currently in alpha. There are several
different ways to run them, some of which are deprecated.

## Plugin Configuration File

Accepting a configuration file as first program argument (legacy
plugins) or in the ResourceList (KRM plugins) is required by the
Kustomize plugin protocol. The configuration file is accepted, but
currently it is not actually used. Note, that if
`yaml-include-transformer` is run with single argument and that
argument is not an option it is regarded as a configuration file.

### Installation as legacy EXEC plugin

A [legacy EXEC
plugins](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/exec_plugins/)
is an executable that accepts a single argument on its command line -
the name of a YAML file containing its configuration (the file name
provided in the kustomization.yaml). The plugin executable must be
located at
`$XDG_CONFIG_HOME/kustomize/plugin/${apiVersion}/LOWERCASE(${kind})/${kind}`. The
default value of `XDG_CONFIG_HOME` is `$HOME/.config`.

To install `yaml-include-transformer` as a legacy EXEC plugin run

```shell
$ yaml-include-transformer --install --legacy --exec
Installing kustomize exec plugin /home/username/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer
copy '/home/username/go/bin/yaml-include-transformer' to '/home/username/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer/YamlIncludeTransformer'
/home/username/go/bin/yaml-include-transformer: Kustomize exec plugin Installation complete
```

Create plugin configuration file in the project directory
(p.e. include-plugin.yaml):

```shell
yaml-include-transformer --plugin-conf -legacy > include-plugin.yaml
```

Add a transformer declaration to the `kustomization.yaml` file:

```yaml
transformers:
  - include-plugin.yaml
```

Invoke kustomize build:

```shell
kustomize build --enable-alpha-plugins 
```
See an example in the `examples/legacy-exec` subdirectory.

### Installation as an Exec KRM function

An 
[Exec KRM function](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/exec_krm_functions/)
is an executable that accepts a ResourceList as input on stdin and
emits a ResourceList as output on stdout. The executable must be
located in the project directory, the exact location is is defined in
the plugin configuration file.

To install `yaml-include-transformer` as an Exec KRM function run in the 
project directory:

```shell
$ yaml-include-transformer --install --krm --exec
```

Create plugin configuration file in the project directory
(p.e. include-plugin.yaml):

```shell
$ yaml-include-transformer --plugin-conf --krm --exec > include-plugin.yaml
```

Add a transformer declaration to the `kustomization.yaml` file:

```yaml
transformers:
  - include-plugin.yaml
```

Invoke kustomize build:

```shell
kustomize build --enable-alpha-plugins --enable-exec
```
See an example in the `examples/krm-exec` subdirectory.

### Installation as Containerized KRM function

A 
[Containerized KRM Function](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/containerized_krm_functions)
is a container whose entrypoint accepts a ResourceList as input on stdin 
and emits a ResourceList as output on stdout.

To use `yaml-include-transformer` as a Containerized KRM function 
create plugin configuration file in the project directory
(p.e. include-plugin.yaml):

```shell
$ yaml-include-transformer --plugin-conf --krm > include-plugin.yaml
```

The plugin configuration contains image tag for the `yaml-include-transformer`
container image ``. 


Add a transformer declaration to the `kustomization.yaml` file:

```yaml
transformers:
  - include-plugin.yaml
```

Invoke kustomize build in the project directory:

```shell
kustomize build --enable-alpha-plugins --mount type=bind,source=".",target=/work
```

This plugin needs to access the project directory so this command mounts the 
project directory into the plugin container.

See an example in the `examples/krm-containerized` subdirectory.



Configuring ArgoCD to use Kustomize with the plugin 
---------------------------------------------------

There are several ways to use `yaml-include-transformer` with
[ArgoCD](https://argoproj.github.io), each one comes with its
advantages and disadvantages.

*WARNING*: Kustomize plugins support is an alpha functionality,
enabling it on your ArgoCD instance may effectively allow anyone with
commit access to the Git repositories to run their code inside your
`argocd-repo-server` pod.

### Using Kustomize legacy EXEC plugin

The `argocd-repo-server` deployment must be customized to to use a
customized docker image that includes the `yaml-include-transformet` binary.
One is also required  to change the `kustomize.buildOptions` 
value in the `argocd-cm` ConfigMap. 

See more in the [ArgoCD
documentation](https://argo-cd.readthedocs.io/en/stable/operator-manual/custom_tools)
on inclusion of custom tools.

#### Building a Customized ArgoCD Image

This command will add the `yaml-include-transformer` binary to the
source ArgoCD docker image and installs it as a customize plugin.  You
can customize target repository and source image using environment
variables, see details in the Makefile.

```shell
$ env ARGOCD_REPO=some-repo/argocd-yit ARGOCD_VER=v2.4.4  make argo_docker_build
echo 	"FROM quay.io/argoproj/argocd:v2.4.4 \n" \
	"ADD ./yaml-include-transformer /usr/local/bin\n" \
	"RUN /usr/local/bin/yaml-include-transformer -i\n" > Dockerfile.argocd
docker build -f Dockerfile.argocd -t some-repo/argocd-yit:v2.4.4_yitv0.0.4alpha1 .
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
Successfully tagged some-repo/argocd-yit:v2.4.4_yitv0.0.4alpha1
```

`make argo_docker_push` will push the image to your repository.

#### Patching the ArgoCD configuration

The following command patches the deployment of `argocd-repo-server` to use the customized
docker image and changes the kustomize command line flags in the parameter `kustomize.buildOptions` 
in the `argocd-cm` ConfigMap. See details in the Makefile.

```shell
 $ /usr/bin/env ARGOCD_REPO=some-repo/argocd-yit ARGOCD_VER=v2.4.4  make argo_patch_legacy_exec
kubectl patch deployment -n  argocd argocd-repo-server -p \
'{"spec" : {"template" : { "spec" : { "containers" : [ { "image" : "some-repo/argocd-yit:v2.4.4_yitv0.0.4alpha1", "name" : "argocd-repo-server"  }]}}}}'
deployment.apps/argocd-repo-server patched
kubectl patch cm -n argocd argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-alpha-plugins"}}'
configmap/argocd-cm patched
```

### Using Exec KRM function

In this mode the binary must be installed inside the repository as
described [above](#installation-as-an-exec-krm-function).  One is also
required to change the `kustomize.buildOptions` value in the
`argocd-cm` ConfigMap:

```shell
$ make argo_patch_krm_exec
kubectl patch cm -n argocd argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-alpha-plugins --enable-exec"}}'
configmap/argocd-cm patched
```

### Using as an ArgoCD CM plugin

One can also use `yaml-include-transformer` as an ArgoCD Configuration
Management Plugin (CMP) without using kustomize. 

There are two ways to set-up CM plugins: using the `argocd-cm` ConfigMap
and using sidecars.

#### Setting up using `argocd-cm`

1. One need to make the binary available in the `argocd-repo-server` container 
   as described [above](#using-kustomize-legacy-EXEC-plugin). 

2. Configure plugin in the `argocd-cm` ConfigMap:

```shell
$ make argo_patch_cmp_cm
[TODO]
```

3. Configure your Application to use the plugin:

```
spec:
  source:
    plugin:
      name: YamlIncludeTransformer
```

#### Setting up using sidecar

[TO-BE-DONE]


Using `kubectl` with the plugin
-------------------------------

Run kustomize, which is built into kubectl.

```shell
kubectl kustomize  --enable-alpha-plugins=true   . 
```

AFAIK currently there is no way to enable plugins when running 
`kubectl apply -k`, but as a workaround one could pipe
kustomize output into kubectl apply command like:

```shell
kubectl kustomize  --enable-alpha-plugins=true . | kubectl apply -f -
```



