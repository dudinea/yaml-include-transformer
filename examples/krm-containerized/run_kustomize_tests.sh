#!/bin/bash
set -e
echo "running $0 in $(pwd)"

BINARY=../../yaml-include-transformer
# ensure it won't run in legacy mode
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"
if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi

echo "testing installation"
"$BINARY" -i --krm  2>&1 | grep  "There is no need to install the binary when using a containerized KRM plugin." 

echo "testing plugin config generation"
"$BINARY" --plugin-conf --krm > plugin-config.yaml.test
diff -u plugin-config.yaml.test plugin-config.yaml

echo "testing yaml processing"
kustomize build --enable-alpha-plugins --mount type=bind,source=".",target=/work > example.out.test
diff -u example.out.test example.out

rm -f example.out.test plugin-config.yaml.test
echo "OK"
