#!/usr/bin/env bash

PS4='\033[34m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "$(test "$(basename "$PWD")" == stage)"
: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  test -f "$STAGE/provisioning/vpn/terraform.tfstate" &&
    jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provisioning/vpn/terraform.tfstate"
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

cd "$STAGE"

digitalocean | while read -r DROPLET; do
  PUBLIC_IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
  PRIVATE_IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address_private')"
  REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
  SERVER_NAME="$REGION.do.vpn.logbook"
  SUBNET_ADDR="$(sed 's;//.*;;g' <"$STAGE/config/digitalocean.jsonc" | jq --arg region "$REGION" -r '.vpn[$region]')"

  if ! test -f "secrets/pki/vpn/issued/$SERVER_NAME.crt"; then
    EASYRSA_PKI="secrets/pki/vpn" easyrsa --batch build-server-full "$SERVER_NAME" nopass
  fi

  scp -i "secrets/ssh/do" \
    "secrets/pki/vpn/ca.crt" \
    "secrets/pki/vpn/issued/$SERVER_NAME.crt" \
    "secrets/pki/vpn/private/$SERVER_NAME.key" \
    "secrets/pki/vpn-users/crl.pem" \
    "secrets/ovpn-auth/ovpn_auth_database.yml" \
    "scripts/deployment/upload/vpn.sh" \
    "$VPS_SUDO_USER@$PUBLIC_IP:/home/$VPS_SUDO_USER/"

  # shellcheck disable=SC2012,SC2087
  ssh -i "secrets/ssh/do" "$VPS_SUDO_USER@$PUBLIC_IP" <<-HERE
    PS4='\033[35m$VPS_SUDO_USER@$PUBLIC_IP \$1:\$LINENO:\033[0m '    
    set -xe

    export VPS_SUDO_USER='$VPS_SUDO_USER'
    export SERVER_NAME='$SERVER_NAME'
    export PRIVATE_IP='$PRIVATE_IP'
    export OPENVPN_SUBNET_ADDRESS='$SUBNET_ADDR'
    export OPENVPN_SUBNET_MASK='255.255.255.0'
    sudo --preserve-env bash vpn.sh
    rm -rfv vpn.sh
HERE
done
