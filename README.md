# kustomize field include processor


A simple YAML processor that implements file include for
YAML files.

## Standalone usage

An example of YAML input:


```yaml
program:
  language: python
  code!textfile:  source.py
  input!base64file: data.bin

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
  language: python
```

## Usage as kustomize plugin

Installation as an "exec" plugin

```shell
./yaml-include-transformer -i
Installing kustomize exec plugin /home/dudin/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer
copy './yaml-include-transformer' to '/home/dudin/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer/FieldIncludeTransformer'
Installation complete
```



Create plugin configuration file (p.e. include-plugin.yaml):

```yaml
apiVersion: kustomize-utils.dudinea.org/v1
kind: FieldIncludeTransformer
metadata:
  name: notImportantHere
```

Add transformer to the `kustomization.yaml` file:

```yaml
transformers:
  - include-plugin.yaml
```
Invoke kustomize build:

```shell
kustomize --enable-exec --enable-alpha-plugins build 
```

## Usage as kustomize shared library based plugin

[TBD]

## Usage as kustomize KRM function based plugin

[TBD]

## Configuring ArgoCD to use the plugin


## Configuration file

Not yet

## Security considerations

[TBD]


## Links

On KRM
  
https://github.com/kubernetes/enhancements/tree/master/keps/sig-cli/2906-kustomize-function-catalog


https://github.com/kubernetes-sigs/kustomize/issues/3936

https://github.com/kubernetes/enhancements/issues/2906

https://github.com/kubernetes/enhancements/tree/master/keps/sig-cli/2953-kustomize-plugin-graduation


https://github.com/kubernetes-sigs/cli-experimental/pull/211

https://github.com/kubernetes/enhancements/issues/2953

https://kubectl.docs.kubernetes.io/guides/extending_kustomize/exec_krm_functions/

https://kubectl.docs.kubernetes.io/guides/extending_kustomize/containerized_krm_functions/


