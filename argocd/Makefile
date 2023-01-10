
ARGOCD_VER?=v2.4.4
YIT_VER?=v0.0.2-alpha
REPO?=quay.io/evgeni_doudine/argocd-yit
ARGO_DOCKERTAG?=$(REPO):a$(ARGOCD_VER)-y$(YIT_VER)

argo_docker_build: yaml-include-transformer
	docker build . -t $(ARGO_DOCKERTAG)

argo_docker_push: argo_docker_build
	docker push $(ARGO_DOCKERTAG)

argo_patch:
	kubectl patch deployment -n argocd argocd-repo-server -p '{"spec" : {"template" : { "spec" : { "containers" : [ { "image" : "$(ARGO_DOCKERTAG)", "name" : "argocd-repo-server"  }]}}}}'
	kubectl patch cm -n argocd argocd-cm -p '{"data" : {"kustomize.buildOptions" : "--enable-exec --enable-alpha-plugins"}}'


