#!/usr/bin/env bash

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  # shellcheck disable=SC2002,SC2046
  cat $(find provisioning -name 'terraform.tfstate') |
    jq -c 'select(.resources | length > 0) | .resources.[] | select(.type == "digitalocean_droplet").instances.[]'
}

# ---------------------------------------------------------------------------- #
# Provisioning
# ---------------------------------------------------------------------------- #

digitalocean | while read -r HOST; do
  echo "$HOST" | jq '.attributes.ipv4_address'

  ssh-keygen -R "$IP"
  ssh-keyscan "$IP" >>~/.ssh/known_hosts # or Private IP
done

# ---------------------------------------------------------------------------- #
# Trust Root CA on MacOS
# ---------------------------------------------------------------------------- #

# security add-trusted-cert -d \
#   -r trustRoot \
#   -k ~/Library/Keychains/login.keychain-db \
#   "${STAGE:?}/secrets/pki/root/ca.crt"
