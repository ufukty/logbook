#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

bash scripts/provision/app-down "$@"
bash scripts/provision/app-up "$@"
bash scripts/deploy ""
