BINARY=yaml-include-transformer
VERSION?=$(shell cat .version)
REPO?=quay.io/evgeni_doudine/yaml-include-transformer
DOCKERTAG?=$(REPO):$(VERSION)

ARGOCD_VER?=v2.4.4
ARGOCD_NS?=argocd
ARGOCD_REPO?=quay.io/evgeni_doudine/argocd-yit
ARGOCD_DOCKERTAG?=$(ARGOCD_REPO):$(ARGOCD_VER)_yit$(VERSION)

$(BINARY): main.go  ## Build the program binary
	go build -ldflags "-X main.version=$(VERSION)"

build_docker:  ## Build docker image for the program 
	docker build -t $(DOCKERTAG) .

push_docker:  
	docker push $(DOCKERTAG) 

install: $(BINARY)
	go install -v

install_plugin: $(BINARY)
	./$(BINARY) -i

clean:
	rm -v -f $(BINARY) examples/example.out Dockerfile.argocd

tests: test_example 

test_example: $(BINARY)
	cd examples && ../$(BINARY) < example.yaml | tee example.out
	diff -u examples/example.out examples/example.test_out 

kustomize_tests: test_install

test_install: $(BINARY)
	cd kustomize-tests/test-install && ./run_tests.sh

argo_docker_build: $(BINARY)
	echo 	"FROM quay.io/argoproj/argocd:$(ARGOCD_VER) \n" \
		"ADD ./$(BINARY) /usr/local/bin\n" \
		"RUN /usr/local/bin/$(BINARY) -i\n" > Dockerfile.argocd
	docker build -f Dockerfile.argocd -t $(ARGOCD_DOCKERTAG) .

argo_docker_push: argo_docker_build
	docker push $(ARGOCD_DOCKERTAG)

argo_patch:
	kubectl patch deployment -n  $(ARGOCD_NS) argocd-repo-server -p \
	'{"spec" : {"template" : { "spec" : { "containers" : [ { "image" : "$(ARGOCD_DOCKERTAG)", "name" : "argocd-repo-server"  }]}}}}'
	kubectl patch cm -n $(ARGOCD_NS) argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-exec --enable-alpha-plugins"}}'

.PHONY: argo_patch argo_docker_push argo_docker_build test_install \
	kustomize_tests test_example tests clean build_docker push_docker install install_plugin 
