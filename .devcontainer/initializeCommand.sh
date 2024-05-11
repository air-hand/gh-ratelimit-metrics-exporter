#!/usr/bin/env bash

HERE=$(dirname $(readlink -f $0))

cd $HERE

cat <<EOF> devcontainer.env
GH_TOKEN=${GH_TOKEN:-$(gh auth token)}
EOF
