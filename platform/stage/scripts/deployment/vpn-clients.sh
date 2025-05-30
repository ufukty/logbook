#!/usr/bin/env bash

PS4='\033[34m$0:$LINENO\033[0m: '
set -xe

: "$(test "$(basename "$PWD")" == stage)"
: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

cd "$STAGE"

test -f "provisioning/vpn/terraform.tfstate" && (
  jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provisioning/vpn/terraform.tfstate" |
    while read -r DROPLET; do
      IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
      REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
      SERVERNAME="$REGION.do.vpn.logbook"

      if ! test -f "secrets/pki/vpn/issued/$SERVERNAME.crt"; then
        EASYRSA_PKI="secrets/pki/vpn" easyrsa --batch build-server-full "$SERVERNAME" nopass
      fi

      scp -i "secrets/ssh/do" \
        "secrets/pki/vpn/ca.crt" \
        "secrets/pki/vpn/issued/$SERVERNAME.crt" \
        "secrets/pki/vpn/private/$SERVERNAME.key" \
        "secrets/pki/vpn-users/crl.pem" \
        "secrets/ovpn-auth/ovpn_auth_database.yml" \
        "$VPS_SUDO_USER@$IP:/home/$VPS_SUDO_USER/"

      # shellcheck disable=SC2012
      # shellcheck disable=SC2087
      ssh -i "secrets/ssh/do" "$VPS_SUDO_USER@$IP" sudo bash <<-HERE
        PS4='\033[33m$VPS_SUDO_USER@$IP \$1:\$LINENO:\033[0m '
        set -xe

        cd "/home/$VPS_SUDO_USER"
        mv ca.crt crl.pem ovpn_auth_database.yml /etc/openvpn/
        mv "$SERVERNAME.crt" /etc/openvpn/server.crt
        mv "$SERVERNAME.key" /etc/openvpn/server.key

        cd /etc/openvpn

        openvpn --genkey secret tls-crypt.key

        chown -R openvpn:openvpn *

        chmod 600 server.key ovpn_auth_database.yml
        chmod 640 ./{ca.crt,server.crt,crl.pem,tls-crypt.key,server.conf}
        chmod 750 ccd
        test $(ls -1 ccd | wc -l) -gt 0 && chmod 640 ccd/*
HERE
    done
)
