#!/usr/bin/env bash

set -xe
PS4='\033[31m$0:$LINENO: \033[0m'

test -d .git || (echo "Run from root folder" && exit 1)

# VERSION="$(date -u +%y%m%d-%H%M%S)-$(git describe --tag --always --dirty)"

mkdir -p build/dev/{linux,darwin}
cd backend
for MAINFILE in cmd/*/main.go; do
  SERVICE="$(basename "$(dirname "$MAINFILE")")"
  for OS in "darwin" "linux"; do
    GOOS=$OS GOARCH=amd64 go build -o "../build/dev/$OS/$SERVICE" "./cmd/$SERVICE"
  done
done
cd ..

tree build/dev
