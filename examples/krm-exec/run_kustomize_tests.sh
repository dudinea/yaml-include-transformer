#!/bin/bash
set -e
echo "running $0 in $(pwd)"
BINARY=../../yaml-include-transformer
# ensure it won't run in legacy mode
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"
if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi
rm -rf plugins/

echo "testing plugin config generation"
"$BINARY" --plugin-conf --krm --exec > plugin-config.yaml.test
diff -u plugin-config.yaml.test plugin-config.yaml

KRMPLUGIN="./plugins/YamlIncludeTransformer"

echo "testing plugin config generation"
"$BINARY" --install --krm --exec
if ! [ -x $KRMPLUGIN  ]; then
    echo "Local install failed to copy plugin to ${KRMPLUGIN}"
    exit 1
fi

echo "testing yaml processing"
kustomize-v4.5.7 build --enable-alpha-plugins --enable-exec  > example.out.test
diff -u example.out.test example.out

rm -f example.out.test plugin-config.yaml.test
echo "OK"
