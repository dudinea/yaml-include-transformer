

build: main.go
	go build

build_docker:
	docker build -t kustomize-field-include:v0.0.1 .

PLUGINDIR=~/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer

install: kustomize-field-include
	./kustomize-field-include -i
	# mkdir -p $(PLUGINDIR)
	# cp kustomize-field-include $(PLUGINDIR)/FieldIncludeTransformer

clean:
	rm -f kustomize-field-include

tests: exec-tests fn-tests

exec-tests:
	cd kustomize-exec-test && kustomize --enable-exec --enable-alpha-plugins build 

fn-tests:
	cd kustomize-exec-test && kustomize --enable-exec --enable-alpha-plugins build 
