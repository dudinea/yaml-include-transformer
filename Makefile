VERSION=$(shell cat .version)

build: main.go
	go build -ldflags "-X main.version=$(VERSION)"

build_docker:
	docker build -t yaml-include-transformer:v0.0.1 .

PLUGINDIR=~/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/yamlincludetransformer

install: yaml-include-transformer
	./yaml-include-transformer -i

clean:
	rm -v -f yaml-include-transformer

test_examples:
	cd examples && ../yaml-include-transformer < example.yaml

kustomize_tests: test_install

#test_multi_yaml test_single_yaml

test_install:
	cd kustomize-tests/test-install && ./run_tests.sh

#exec-tests:
#	cd kustomize-tests/test-install && kustomize --enable-exec --enable-alpha-plugins build 

