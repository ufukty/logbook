#!/usr/bin/env bash

set -xe

: "${EASYRSA_SERVER_NAME:?}"

# ---------------------------------------------------------------------------- #
# Certificates
# ---------------------------------------------------------------------------- #

cd /etc/openvpn/easy-rsa

EASYRSA_CA_NAME="$SERVER_NAME-certificate-authority"
EASYRSA_SERVER_NAME="$SERVER_NAME-server"

# Create the PKI, set up the CA, the DH params and the server certificate
./easyrsa --vars=./vars init-pki
./easyrsa --vars=./vars --batch --req-cn="${EASYRSA_CA_NAME:?}" build-ca nopass

./easyrsa --vars=./vars --batch build-server-full "$EASYRSA_SERVER_NAME" nopass
EASYRSA_CRL_DAYS=3650 ./easyrsa --vars=./vars gen-crl

openvpn --genkey secret /etc/openvpn/tls-crypt.key

# ---------------------------------------------------------------------------- #
# Upload Ovpn-Auth database
# ---------------------------------------------------------------------------- #

: "${REMOTE_USER:?}"

test -f provision/vpn/terraform.tfstate && (
  # shellcheck disable=SC2002
  cat provision/vpn/terraform.tfstate |
    jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' |
    while read -r DROPLET; do
      IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
      scp -i secrets/ssh/do \
        secrets/ovpn-auth/ovpn_auth_database.yml \
        "$REMOTE_USER@$IP:ovpn_auth_database.yml"
      ssh -i secrets/ssh/do -n "$REMOTE_USER@$IP" "sudo bash -c '
        set -xe
        mv /home/$REMOTE_USER/ovpn_auth_database.yml /etc/openvpn/
        chmod 744 /etc/openvpn/ovpn_auth_database.yml
        chown root:root /etc/openvpn/ovpn_auth_database.yml
      '"
    done
)
