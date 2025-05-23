#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

# ---------------------------------------------------------------------------- #
# Environment vars
# ---------------------------------------------------------------------------- #

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# Declarations
# ---------------------------------------------------------------------------- #

terraform-apply() (
  cd "$1"
  terraform apply --auto-approve --var-file="$STAGE/provisioning/vars.tfvars"
)

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

terraform-apply "$STAGE/provisioning/vpc"
terraform-apply "$STAGE/provisioning/vpn"

bash "$STAGE/scripts/provision/update-ssh-conf.sh"
echo "Connect vpn in separate tab [Enter]" && read -r
bash "$STAGE/scripts/provision/ssh-keyscan.sh"
bash "$STAGE/scripts/provision/update-dns-records.sh"

terraform-apply "$STAGE/provisioning/application"

bash "$STAGE/scripts/provision/update-ssh-conf.sh"
bash "$STAGE/scripts/provision/ssh-keyscan.sh"
bash "$STAGE/scripts/provision/update-dns-records.sh"
