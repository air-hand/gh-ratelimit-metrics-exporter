FROM homebrew/brew:latest as brew

FROM prom/prometheus:latest as prometheus

FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm as builder

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ENV \
    DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    GO111MODULE=on

COPY --from=brew /home/linuxbrew /home/linuxbrew

COPY --from=prometheus /bin/promtool /usr/local/bin/promtool

RUN apt-get update && \
    apt-get install -y curl
#    ./aqua-installer -v v2.28.0 && \
##    rm aqua-installer && \
#    apt-get remove -y curl && \
#    apt-get clean

WORKDIR /work

RUN chown -R vscode:vscode /work \
    && chown -R vscode:vscode /go

USER vscode

ENV HOME /home/vscode

RUN curl -sSfL -o /tmp/aqua-installer https://raw.githubusercontent.com/aquaproj/aqua-installer/v3.0.1/aqua-installer && \
    echo "fb4b3b7d026e5aba1fc478c268e8fbd653e01404c8a8c6284fdba88ae62eda6a  /tmp/aqua-installer" | sha256sum -c && \
    chmod +x /tmp/aqua-installer && \
    /tmp/aqua-installer -v v2.28.0 && \
    rm /tmp/aqua-installer && \
    go install github.com/bitnami/wait-for-port@v1.0.7 \
    ;

ENV PATH "${AQUQ_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/aquaproj-aqua}/bin:${PATH}"

COPY --chown=vscode:vscode aqua.yaml Makefile go.mod go.sum ./
COPY --chown=vscode:vscode ./app ./app

RUN aqua i
