#!/bin/bash
set -e
echo "Running install-test"
#set -x
# FIXME: write real tests that parse and normalize yaml

BINARY=../../yaml-include-transformer
PLUGINDIR="${HOME}/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer"

if [ -d "$PLUGINDIR" ]; then
    rm -r -f "$PLUGINDIR"
fi

"$BINARY" -p > plugin.yaml
"$BINARY" -i 2>&1 | grep "Kustomize exec plugin Installation complete" > /dev/null
RESULT=$(kustomize build --enable-exec --enable-alpha-plugins)
EXPECTED=$(cat <<EOF
apiVersion: v1
kind: fooo
metadata:
  name: someobj
text: |
  Once upon a time
  There was a text processing tool
  That could not include files.
EOF
)

if [[ $RESULT != $EXPECTED ]]; then
    echo -e "Unexpected result:\n"
    echo "$RESULT" 
    echo -e "Expected result:\n"
    echo "$EXPECTED"
    exit 2
fi
echo "OK"
