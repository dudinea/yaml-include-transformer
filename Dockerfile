FROM stakater/base-alpine:latest
ADD ./kustomize-field-include /usr/local/bin/
ENTRYPOINT /usr/local/bin/kustomize-field-include

