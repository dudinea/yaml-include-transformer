#!/bin/bash
set -e
echo "Running install-test"

BINARY=../../yaml-include-transformer
# ensure it won't run in legacy mode
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"
if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi

# FIXME
#"$BINARY" -p > plugin-config.yaml.test
#diff -u plugin-config.yaml.test plugin-config.yaml

kustomize build --enable-alpha-plugins --mount type=bind,source=".",target=/work > example.out.test
diff -u example.out.test example.out

rm -f example.out.test plugin-config.yaml.test
echo "OK"
