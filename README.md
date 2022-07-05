# YAML field include processor


A simple YAML processor that implements file include for
YAML files.

An example of YAML input:


```yaml
program:
  code!textfile:  source.py
  language: python
```

Run yaml processor:

```shell
yaml-field-include < test.yaml
```

Output:

```yaml
program:
  code: |
    print("Hello!\n")
  language: python
```

## Usage as kustomize plugin



## Configuration file


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


