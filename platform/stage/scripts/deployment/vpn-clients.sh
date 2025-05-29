#!/usr/bin/env bash

PS4='\033[34m$0:$LINENO\033[0m: '
set -xe

: "$(test "$(basename "$PWD")" == stage)"
: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

test -f "$STAGE/provisioning/vpn/terraform.tfstate" && (
  jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provisioning/vpn/terraform.tfstate" |
    while read -r DROPLET; do
      IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
      REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
      SERVERNAME="$REGION.do.vpn.logbook"

      if ! test -f "$STAGE/secrets/pki/vpn/issued/$SERVERNAME.crt"; then
        EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa --batch build-server-full "$SERVERNAME" nopass
      fi

      scp -i "$STAGE/secrets/ssh/do" \
        "$STAGE/secrets/pki/vpn/ca.crt" \
        "$STAGE/secrets/pki/vpn/issued/$SERVERNAME.crt" \
        "$STAGE/secrets/pki/vpn/private/$SERVERNAME.key" \
        "$STAGE/secrets/pki/vpn-users/crl.pem" \
        "$STAGE/secrets/ovpn-auth/ovpn_auth_database.yml" \
        "$VPS_SUDO_USER@$IP:/home/$VPS_SUDO_USER"

      # shellcheck disable=SC2087
      ssh -i "$STAGE/secrets/ssh/do" "$VPS_SUDO_USER@$IP" sudo bash <<-HERE
        PS4='\033[33m$VPS_SUDO_USER@$IP \$LINENO:\033[0m '
        set -xe
        
        cd /home/$VPS_SUDO_USER
        mv ca.crt crl.pem ovpn_auth_database.yml /etc/openvpn/
        mv "$SERVERNAME.crt" /etc/openvpn/server.crt
        mv "$SERVERNAME.key" /etc/openvpn/server.key

        cd /etc/openvpn
        chmod 644 crl.pem
        chmod 744 ovpn_auth_database.yml
        chown -R root:root *
HERE
    done
)
