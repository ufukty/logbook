#!/usr/local/bin/bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

cd "${STAGE:?}/provisioning/vpn"
terraform apply --auto-approve --var-file="${STAGE:?}/provisioning/vars.tfvars"

bash scripts/provision/aggregate-ssh-conf

echo "Connect vpn in separate tab [Enter]"
# shellcheck disable=SC2162
read # wait

bash scripts/digitalocean/ssh-keyscan
bash scripts/provision/update-dns-records
