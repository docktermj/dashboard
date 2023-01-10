#!/bin/bash

# Reference: https://stackoverflow.com/questions/54933983/how-to-serve-documentation-using-godoc-together-with-go-modules

set -x  # optional

REPO_HOME=$(echo ${PWD})
PORT=6060

docker run \
    --env "GOPATH=/tmp/go" \
    --name godoc \
    --publish ${PORT}:80 \
    --rm \
    --volume ${REPO_HOME}:/tmp/go/src/ \
    golang \
        bash -c "go get golang.org/x/tools/cmd/godoc \
            && echo http://localhost:80/pkg/ \
            && /tmp/go/bin/godoc -http=:80"
