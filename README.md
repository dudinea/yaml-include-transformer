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


