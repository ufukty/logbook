#!/usr/local/bin/bash

cd "${STAGE:?}/provisioning/application"
terraform destroy "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"

bash scripts/provision/aggregate-ssh-conf
