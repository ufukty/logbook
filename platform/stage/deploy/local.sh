#!/usr/bin/env bash

PS4='\033[34m$0:$LINENO\033[0m:'
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  # shellcheck disable=SC2002,SC2046
  cat $(find "$STAGE/provision" -name 'terraform.tfstate') |
    jq -c 'select(.resources | length > 0) | .resources.[] | select(.type == "digitalocean_droplet").instances.[]'
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

digitalocean | while read -r HOST; do
  IP="$(echo "$HOST" | jq -r '.attributes.ipv4_address')"
  ssh-keygen -R "$IP"
  ssh-keyscan "$IP" >>~/.ssh/known_hosts # or Private IP
done

# ---------------------------------------------------------------------------- #
# Trust Root CA on MacOS
# ---------------------------------------------------------------------------- #

CHAIN="$HOME/Library/Keychains/login.keychain-db"

if security find-certificate -c "Logbook Stage Root CA" "$CHAIN" >/dev/null 2>&1; then
  security delete-certificate -c "Logbook Stage Root CA" "$CHAIN"
fi

security add-trusted-cert -d -r trustRoot -k "$CHAIN" "$STAGE/secrets/pki/root/ca.crt"
