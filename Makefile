

build: main.go
	go build

build_docker:
	docker build -t yaml-include-transformer:v0.0.1 .

PLUGINDIR=~/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer

install: yaml-include-transformer
	./yaml-include-transformer -i
	# mkdir -p $(PLUGINDIR)
	# cp yaml-include-transformer $(PLUGINDIR)/FieldIncludeTransformer

clean:
	rm -f yaml-include-transformer

tests: exec-tests fn-tests

exec-tests:
	cd kustomize-exec-test && kustomize --enable-exec --enable-alpha-plugins build 

