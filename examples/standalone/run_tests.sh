#!/bin/bash
set -e
echo "running $0 in $(pwd)"

BINARY=../../yaml-include-transformer

echo testing example-multidoc.yaml
"${BINARY}" < example.yaml  > example.out.test
diff -u example.out.test example.out

echo testing example.yaml
"${BINARY}" < example-multidoc.yaml  > example-multidoc.out.test
diff -u example-multidoc.out.test example-multidoc.out


rm -f example.out.test example-multidoc.out.test
echo "OK"
