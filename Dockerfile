#FROM homebrew/brew:latest as brew
#
#FROM prom/prometheus:latest as prometheus
#
#FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm as builder
#
#SHELL ["/bin/bash", "-o", "pipefail", "-c"]
#
#ENV \
#    DEBIAN_FRONTEND=noninteractive \
#    LANG=en_US.UTF-8 \
#    GO111MODULE=on
#
#COPY --from=brew /home/linuxbrew /home/linuxbrew
#
#COPY --from=prometheus /bin/promtool /usr/local/bin/promtool
#
#WORKDIR /work
#
#RUN chown -R vscode:vscode /work \
#    && chown -R vscode:vscode /go
#
#RUN apt-get update && \
#    apt-get install -y curl && \
#    curl -sSfL -O https://raw.githubusercontent.com/aquaproj/aqua-installer/v3.0.1/aqua-installer && \
#    echo "fb4b3b7d026e5aba1fc478c268e8fbd653e01404c8a8c6284fdba88ae62eda6a  aqua-installer" | sha256sum -c && \
#    chmod +x aqua-installer && \
#    ./aqua-installer -v v2.28.0 && \
#    rm aqua-installer && \
#    apt-get remove -y curl && \
#    apt-get clean
#
#USER vscode
#
#RUN \
#    go install github.com/bitnami/wait-for-port@v1.0.7 \
#    ;
#
#COPY --chown=vscode:vscode aqua.yaml Makefile go.mod go.sum ./
#COPY --chown=vscode:vscode ./app ./app
#
#RUN aqua i && make build
#
FROM gcr.io/distroless/base-debian12:latest as release

COPY --chown=nonroot:nonroot gh-ratelimit-metrics-exporter /app

USER nonroot

ENTRYPOINT ["/app"]
