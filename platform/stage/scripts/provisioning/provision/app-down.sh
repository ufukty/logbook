#!/usr/local/bin/bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

cd "${STAGE:?}/provisioning/application"
terraform destroy "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"

bash scripts/provision/aggregate-ssh-conf
