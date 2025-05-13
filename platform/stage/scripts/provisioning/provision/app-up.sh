#!/usr/local/bin/bash

set -xe

cd "${STAGE:?}/provisioning/application"
terraform apply "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"

bash scripts/provision/aggregate-ssh-conf
bash scripts/provision/ssh-keyscan
bash scripts/provision/update-dns-records
