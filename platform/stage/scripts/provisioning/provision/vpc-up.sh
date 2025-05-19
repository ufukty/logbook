#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

cd "${STAGE:?}/provisioning/vpc"
terraform apply "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
