#!/bin/bash
# Si Unix (Linux / MacOS)
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=$os -e GOARCH=$arch golang:1.22 go build -v