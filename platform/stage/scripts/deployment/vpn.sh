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
  SUBNET_ADDR="$(sed 's;//.*;;g' <"$STAGE/config/digitalocean.jsonc" | jq --arg region fra1 -r '.vpn[$region]')"

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
  ssh -i "secrets/ssh/do" "$VPS_SUDO_USER@$PUBLIC_IP" sudo bash <<-HERE
    PS4='\033[33m$VPS_SUDO_USER@$PUBLIC_IP \$1:\$LINENO:\033[0m '
    set -xe

    export USER_ACCOUNT_NAME='$VPS_SUDO_USER'
    export SERVER_NAME='$SERVER_NAME'
    export PUBLIC_IP='$PUBLIC_IP'
    export PRIVATE_IP='$PRIVATE_IP'
    export OPENVPN_SUBNET_ADDRESS='$SUBNET_ADDR'
    export OPENVPN_SUBNET_MASK='255.255.255.0'
    export PUBLIC_ETHERNET_INTERFACE='eth0'
    export PRIVATE_ETHERNET_INTERFACE='eth1'
    
    cd ~
    sudo --preserve-env bash deployment.sh

    cd "/home/$VPS_SUDO_USER"
    mv ca.crt crl.pem ovpn_auth_database.yml /etc/openvpn/
    mv "$SERVER_NAME.crt" /etc/openvpn/server.crt
    mv "$SERVER_NAME.key" /etc/openvpn/server.key

    cd /etc/openvpn

    openvpn --genkey secret tls-crypt.key

    chown -R openvpn:openvpn *

    chmod 600 server.key ovpn_auth_database.yml
    chmod 640 ./{ca.crt,server.crt,crl.pem,tls-crypt.key,server.conf}
    chmod 750 ccd
    test \$(ls -1 ccd | wc -l) -gt 0 && chmod 640 ccd/*

    # sudo bash <<PRIVILEGED
    #   systemctl restart systemd-journald
    #   systemctl restart iptables-activation
    #   sed -E -in-place \"s;$VPS_SUDO_USER(.*)NOPASSWD:(.*);$VPS_SUDO_USER \1 \2;\" /etc/sudoers  
    # PRIVILEGED
HERE
done
