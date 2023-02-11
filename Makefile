BINARY=yaml-include-transformer
VERSION?=$(shell cat .version)
REPO?=quay.io/evgeni_doudine/yaml-include-transformer
DOCKERTAG?=$(REPO):$(VERSION)

ARGOCD_VER?=v2.4.4
ARGOCD_NS?=argocd
ARGOCD_SRC_REPO?=quay.io/argoproj/argocd
ARGOCD_REPO?=quay.io/evgeni_doudine/argocd-yit
ARGOCD_DOCKERTAG?=$(ARGOCD_REPO):$(ARGOCD_VER)_yit$(VERSION)

LDFLAGS=-X main.version=$(VERSION) -X main.dockertag=$(DOCKERTAG)

$(BINARY): $(wildcard pkg/**/*.go) main.go go.mod
	go build -ldflags "$(LDFLAGS)"

build_docker: $(BINARY)
	docker build -t $(DOCKERTAG) .

push_docker:  
	docker push $(DOCKERTAG) 

install: $(BINARY)
	go install -v -ldflags "$(LDFLAGS)"

install_plugin: $(BINARY)
	./$(BINARY) -i

clean:
	rm -v -f $(BINARY) examples/*/example.out.test Dockerfile.argocd

tests: test_standalone kustomize_tests

test_standalone: $(BINARY)
	cd examples/standalone && ./run_tests.sh

kustomize_tests: test_legacy_exec test_krm_containerized test_krm_exec

test_legacy_exec: $(BINARY)
	cd examples/legacy-exec && ./run_kustomize_tests.sh

test_krm_containerized: $(BINARY)
	cd examples/krm-containerized && ./run_kustomize_tests.sh

test_krm_exec: $(BINARY)
	cd examples/krm-exec && ./run_kustomize_tests.sh


argo_docker_build: $(BINARY)
	echo 	"FROM $(ARGOCD_SRC_REPO):$(ARGOCD_VER) \n" \
		"ADD ./$(BINARY) /usr/local/bin\n" \
		"RUN /usr/local/bin/$(BINARY) -i\n" > Dockerfile.argocd
	docker build -f Dockerfile.argocd -t $(ARGOCD_DOCKERTAG) .

argo_docker_push: argo_docker_build
	docker push $(ARGOCD_DOCKERTAG)

argo_patch_image:
	kubectl patch deployment -n  $(ARGOCD_NS) argocd-repo-server -p \
	'{"spec" : {"template" : { "spec" : { "containers" : [ { "image" : "$(ARGOCD_DOCKERTAG)", "name" : "argocd-repo-server"  }]}}}}'

argo_patch_legacy_exec: argo_patch_image
	kubectl patch cm -n $(ARGOCD_NS) argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-alpha-plugins"}}'

argo_patch_krm_exec:
	kubectl patch cm -n $(ARGOCD_NS) argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-alpha-plugins --enable-exec"}}'

argo_patch_cmp_cm: argo_patch_image
	kubectl patch cm -n $(ARGOCD_NS) argocd-cm -p '{"data" : {"configManagementPlugins": "[ { \"name\":  \"YamlIncludeTransformer\", \"generate\": { \"command\" : [ \"/usr/local/bin/yaml-include-transformer\" ],  \"args\": [ \"-f\" , \".\" ]}}]"}}'


.PHONY: argo_patch_legacy_exec argo_docker_push argo_docker_build \
	argo_patch_krm_exec argo_patch_cmp_cm \
	kustomize_tests test_standalone tests \
	test_legacy_exec test_krm_containerized test_krm_exec \
	clean build_docker push_docker install install_plugin 
