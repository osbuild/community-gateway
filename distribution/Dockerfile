FROM registry.redhat.io/ubi9/ubi-minimal

RUN microdnf --nodocs upgrade -y
RUN rpm --import 'https://rpm.dl.getenvoy.io/public/gpg.CF716AF503183491.key'
RUN curl -sL 'https://rpm.dl.getenvoy.io/public/config.rpm.txt?distro=el&codename=9' > /etc/yum.repos.d/envoy.repo
RUN microdnf install -y getenvoy-envoy
RUN microdnf clean all

EXPOSE 10000
CMD "/usr/bin/envoy"
