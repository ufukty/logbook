#!/usr/local/bin/bash

PROGRAM="${1:?}"

set -xe
PS4='\033[31m$0:$LINENO: \033[0m'

cd platform/stage
./commands deploy "$PROGRAM" "$@"
