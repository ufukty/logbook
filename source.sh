#!/usr/local/bin/bash

export WORKSPACE="$PWD"

# shellcheck disable=SC1090
for f in *.env; do . "$f"; done
