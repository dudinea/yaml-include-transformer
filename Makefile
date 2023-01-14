VERSION=$(shell cat .version)
BINARY=yaml-include-transformer

$(BINARY): main.go
	go build -ldflags "-X main.version=$(VERSION)"

build_docker:
	docker build -t yaml-include-transformer:$(VERSION) .

install: $(BINARY)
	go install -v

install_plugin: $(BINARY)
	./$(BINARY) -i

clean:
	rm -v -f $(BINARY) examples/example.out

tests: test_example kustomize_tests

test_example: $(BINARY)
	cd examples && ../$(BINARY) < example.yaml | tee example.out
	diff -u examples/example.out examples/example.test_out 

kustomize_tests: test_install

test_install:
	cd kustomize-tests/test-install && ./run_tests.sh



