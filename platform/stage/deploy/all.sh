#!/usr/bin/env bash

PS4='\033[34m$0:$LINENO\033[0m:'
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  # shellcheck disable=SC2002,SC2046
  cat $(find find "$STAGE/provision" -name 'terraform.tfstate') |
    jq -c 'select(.resources | length > 0) | .resources.[] | select(.type == "digitalocean_droplet").instances.[]'
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

digitalocean | while read -r HOST; do
  IP="$(echo "$HOST" | jq -r '.attributes.ipv4_address')"

  scp -i secrets/ssh/do secrets/pki/root/ca.crt "$VPS_SUDO_USER@$IP:/home/$VPS_SUDO_USER/"

  # shellcheck disable=SC2087
  ssh -i secrets/ssh/do "$VPS_SUDO_USER@$IP" sudo bash <<-HERE
    PS4='\033[33m$VPS_SUDO_USER@$IP \$1:\$LINENO:\033[0m '
    set -xe
    
    cd /home/$VPS_SUDO_USER

    mv ca.crt /usr/local/share/ca-certificates/logbook-stage-root-ca.crt
    update-ca-certificates
    openssl x509 -in /usr/local/share/ca-certificates/logbook-stage-root-ca.crt -noout -text | grep "Logbook Stage Root CA"
HERE
done
