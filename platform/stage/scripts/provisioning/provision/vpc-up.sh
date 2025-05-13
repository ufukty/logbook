#!/usr/local/bin/bash

cd "${STAGE:?}/provisioning/vpc"
terraform apply "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
