FROM quay.io/argoproj/argocd:v2.4.4
ADD ./yaml-include-transformer /usr/local/bin/
RUN /usr/local/bin/yaml-include-transformer -i
