FROM quay.io/argoproj/argocd:latest
ADD ./yaml-include-transformer /usr/local/bin/
RUN /usr/local/bin/yaml-include-transformer -i
