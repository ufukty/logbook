#!/usr/bin/env bash

PS4='\033[32m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${CLIENT_NAME:?}"
: "${EASYRSA_SERVER_NAME:?}"
: "${PUBLIC_IP:?}"
: "${USER_ACCOUNT_NAME:?}"

# ---------------------------------------------------------------------------- #
# Client certs
# ---------------------------------------------------------------------------- #

export DEBIAN_FRONTEND=noninteractive

cd /etc/openvpn/easy-rsa
./easyrsa build-client-full "$CLIENT_NAME" nopass

# ---------------------------------------------------------------------------- #
# Templating
# ---------------------------------------------------------------------------- #

EASYRSA_CA_CERT_CONTENT="$(cat "/etc/openvpn/easy-rsa/pki/ca.crt")"
EASYRSA_CLIENT_KEY_CONTENT="$(cat "/etc/openvpn/easy-rsa/pki/private/$CLIENT_NAME.key")"
EASYRSA_CLIENT_CERT_CONTENT="$(awk '/BEGIN/,/END/' "/etc/openvpn/easy-rsa/pki/issued/$CLIENT_NAME.crt")"
TLS_SIG_KEY_CONTENT="$(cat "/etc/openvpn/tls-crypt.key")"

mkdir -p "/home/$USER_ACCOUNT_NAME/artifacts"

sed \
  -e "s;{{PUBLIC_IP}};$PUBLIC_IP;g" \
  -e "s;{{EASYRSA_CA_CERT_CONTENT}};$EASYRSA_CA_CERT_CONTENT;g" \
  -e "s;{{EASYRSA_CLIENT_CERT_CONTENT}};$EASYRSA_CLIENT_CERT_CONTENT;g" \
  -e "s;{{EASYRSA_CLIENT_KEY_CONTENT}};$EASYRSA_CLIENT_KEY_CONTENT;g" \
  -e "s;{{TLS_SIG_KEY_CONTENT}};$TLS_SIG_KEY_CONTENT;g" \
  -e "s;{{EASYRSA_SERVER_NAME}};$EASYRSA_SERVER_NAME;g" \
  </etc/openvpn/client.ovpn.tpl \
  >"/home/$USER_ACCOUNT_NAME/artifacts/$CLIENT_NAME.ovpn"

chown -R "$USER_ACCOUNT_NAME" "/home/$USER_ACCOUNT_NAME"
