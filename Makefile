

build: main.go
	go build

build_docker:
	docker build -t kustomize-field-include:v0.0.1 .

PLUGINDIR=~/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer

install: kustomize-field-include
	mkdir -p $(PLUGINDIR)
	cp kustomize-field-include $(PLUGINDIR)/FieldIncludeTransformer

clean:
	rm -f kustomize-field-include

