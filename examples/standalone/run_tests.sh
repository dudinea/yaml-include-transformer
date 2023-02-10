#!/bin/bash
set -e
echo "running $0 in $(pwd)"

BINARY=../../yaml-include-transformer

echo testing example.yaml as stdin
"${BINARY}" < example.yaml  > example.out.test
diff -u example.out.test example.out

echo testing example.yaml 
"${BINARY}"  -f example.yaml  > example.out.test
diff -u example.out.test example.out

echo testing example-multidoc.yaml as stdin
"${BINARY}" < example-multidoc.yaml  > example-multidoc.out.test
diff -u example-multidoc.out.test example-multidoc.out

echo testing example-multidoc.yaml
"${BINARY}" -f example-multidoc.yaml  > example-multidoc.out.test
diff -u example-multidoc.out.test example-multidoc.out

echo testing multiple input files
"${BINARY}" -f example-multidoc.yaml example.yaml  > example-multiple.out.test
diff -u example-multiple.out.test example-multiple.out

echo testing directory as input
"${BINARY}" -P  '^example.*\.yaml$' -d -f .  > example-dir.out.test
diff -u example-dir.out.test example-dir.out

rm -f example.out.test example-multidoc.out.test example-multiple.out.test example-dir.out.test
echo "OK"
