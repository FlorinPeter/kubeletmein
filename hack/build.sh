#!/bin/sh

OUTPUT="${OUTPUT:-kubeletmein}"
export CGO_ENABLED=0

OSTYPE=${OSTYPE:-linux}
if [[ ${OSTYPE:0:5} == linux ]]; then
    echo "build for linux"
    export GOOS='darwin'

elif [[ ${OSTYPE:0:6} == darwin ]]; then
    echo "build for mac"
    export GOOS='darwin'
    export GOARCH='amd64'

else
    echo "Unknown os"
fi

time go build -o "${OUTPUT}" ./cmd/kubeletmein
