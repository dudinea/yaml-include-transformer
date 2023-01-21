#!/bin/bash
set -e
echo "running $0 in $(pwd)"
BINARY=../../yaml-include-transformer
# ensure it won't run in legacy mode
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"
if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi

# FIXME
#"$BINARY" -p > plugin-config.yaml.test
#diff -u plugin-config.yaml.test plugin-config.yaml

mkdir -p plugins
cp "${BINARY}" plugins/

kustomize-v4.5.7 build --enable-alpha-plugins --enable-exec  > example.out.test
diff -u example.out.test example.out

rm -f example.out.test plugin-config.yaml.test
echo "OK"
