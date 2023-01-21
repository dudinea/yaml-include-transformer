FROM stakater/base-alpine:latest
ADD ./yaml-include-transformer /usr/local/bin/
RUN mkdir -p /work
ENTRYPOINT /usr/local/bin/yaml-include-transformer
WORKDIR /work

