#!/usr/bin/env bash

# Based on: https://gist.github.com/QueuingKoala/e2c1c067a312384915b5

PS4='\033[33m$0:$LINENO\033[0m '
set -xe

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# PKI dir
# ---------------------------------------------------------------------------- #

rm -rfv "$STAGE/secrets/pki"
mkdir -p "$STAGE"/secrets/pki

# ---------------------------------------------------------------------------- #
# Root CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch --req-cn="Logbook Stage Root CA" build-ca nopass

# ---------------------------------------------------------------------------- #
# Web Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/web" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/web" easyrsa --batch --req-cn="Logbook Stage Web CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch import-req "$STAGE/secrets/pki/web/reqs/ca.req" web
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch sign-req ca web

cp "$STAGE/secrets/pki/root/issued/web.crt" \
  "$STAGE/secrets/pki/web/ca.crt"

# ---------------------------------------------------------------------------- #
# Vpn Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa --batch --req-cn="Logbook Stage VPN CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch import-req "$STAGE/secrets/pki/vpn/reqs/ca.req" vpn
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch sign-req ca vpn

cp "$STAGE/secrets/pki/root/issued/vpn.crt" \
  "$STAGE/secrets/pki/vpn/ca.crt"

# ---------------------------------------------------------------------------- #
# Vpn Users Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" easyrsa --batch --req-cn="Logbook Stage VPN Users CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch import-req "$STAGE/secrets/pki/vpn-users/reqs/ca.req" vpn-users
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --batch sign-req ca vpn-users

cp "$STAGE/secrets/pki/root/issued/vpn-users.crt" \
  "$STAGE/secrets/pki/vpn-users/ca.crt"

EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" EASYRSA_CRL_DAYS=3650 easyrsa --batch gen-crl

# ---------------------------------------------------------------------------- #
# Trust Root CA on MacOS
# ---------------------------------------------------------------------------- #

security add-trusted-cert -d \
  -r trustRoot \
  -k ~/Library/Keychains/login.keychain-db \
  "${STAGE:?}/secrets/pki/root/ca.crt"
