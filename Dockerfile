FROM stakater/base-alpine:latest
ADD ./yaml-include-transformer /usr/local/bin/
ENTRYPOINT /usr/local/bin/yaml-include-transformer

