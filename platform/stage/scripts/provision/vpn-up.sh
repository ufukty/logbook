#!/usr/local/bin/bash

set -xe

cd "${STAGE:?}/provisioning/vpn"
terraform apply --auto-approve --var-file="${STAGE:?}/provisioning/vars.tfvars"

bash scripts/provision/aggregate-ssh-conf

echo "Connect vpn in separate tab [Enter]"
read # wait

bash scripts/digitalocean/ssh-keyscan
bash scripts/provision/update-dns-records
