

build: main.go
	go build

PLUGINDIR=~/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer

install: kustomize-field-include
	mkdir -p $(PLUGINDIR)
	cp kustomize-field-include $(PLUGINDIR)/FieldIncludeTransformer

clean:
	rm -f kustomize-field-include

