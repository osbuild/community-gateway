FROM quay.io/fedora/fedora:latest AS builder
WORKDIR /oidc-authorizer
RUN dnf install -y golang
COPY oidc-authorizer .
ENV GOFLAGS=-mod=vendor
RUN go build -o oidc-authorizer /oidc-authorizer/cmd/oidc-authorizer

FROM quay.io/fedora/fedora-minimal:latest
COPY --from=builder /oidc-authorizer/oidc-authorizer /oidc-authorizer
