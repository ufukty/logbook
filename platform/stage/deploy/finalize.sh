#!/usr/bin/env bash
#
# Performs the final provision on hosts that makes their further unattended
# configutation much harder or impossible due to removal of passwordless
# sudo rights.

PS4='\033[34m$0:$LINENO\033[0m:'
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${VPS_SUDO_USER:?}"

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

  # shellcheck disable=SC2087
  ssh -i secrets/ssh/do "$VPS_SUDO_USER:$IP" sudo bash <<HERE
    systemctl restart systemd-journald
    systemctl restart iptables-activation
    sed -E -in-place \"s;$VPS_SUDO_USER(.*)NOPASSWD:(.*);$VPS_SUDO_USER \1 \2;\" /etc/sudoers
HERE
done
