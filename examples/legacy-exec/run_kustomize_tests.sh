#!/bin/bash
set -e
echo "running $0 in $(pwd)"

BINARY=../../yaml-include-transformer
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"

if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi

"$BINARY" -i 2>&1 | grep "Kustomize exec plugin Installation complete"

"$BINARY" -p > plugin-config.yaml.test
diff -u plugin-config.yaml.test plugin-config.yaml

kustomize build --enable-alpha-plugins > example.out.test
diff -u example.out.test example.out

rm -f example.out.test plugin-config.yaml.test
echo "OK"
