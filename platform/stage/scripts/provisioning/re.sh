#!/usr/local/bin/bash

set -x

bash scripts/provision/app-down "$@"
bash scripts/provision/app-up "$@"
bash scripts/deploy ""
