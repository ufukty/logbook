#!/usr/bin/env bash

PS4='\033[33m$0:$LINENO\033[0m '
set -xe

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# PKI dir
# ---------------------------------------------------------------------------- #

rm -rfv "$STAGE/secrets/pki"
mkdir -p "$STAGE"/secrets/pki/{root,web,vpn}

# ---------------------------------------------------------------------------- #
# Root CA
# ---------------------------------------------------------------------------- #

cd "$STAGE/secrets/pki/root"
easyrsa init-pki
easyrsa --batch --req-cn="Logbook Stage CA" build-ca nopass

# ---------------------------------------------------------------------------- #
# Web Intermediate CA
# ---------------------------------------------------------------------------- #

cd "$STAGE/secrets/pki/web"
easyrsa init-pki
easyrsa --batch gen-req web nopass

cd "$STAGE/secrets/pki/root"
easyrsa --batch import-req "$STAGE/secrets/pki/web/pki/reqs/web.req" web
easyrsa --batch sign-req ca web

cp "$STAGE/secrets/pki/root/pki/issued/web.crt" \
  "$STAGE/secrets/pki/web/pki/ca.crt"

# ---------------------------------------------------------------------------- #
# Vpn Intermediate CA
# ---------------------------------------------------------------------------- #

cd "$STAGE/secrets/pki/vpn"
easyrsa init-pki
easyrsa --batch gen-req vpn nopass

cd "$STAGE/secrets/pki/root"
easyrsa --batch import-req "$STAGE/secrets/pki/vpn/pki/reqs/vpn.req" vpn
easyrsa --batch sign-req ca vpn

cp "$STAGE/secrets/pki/root/pki/issued/vpn.crt" \
  "$STAGE/secrets/pki/vpn/pki/ca.crt"

# ---------------------------------------------------------------------------- #
# Trust Root CA on MacOS
# ---------------------------------------------------------------------------- #

security add-trusted-cert -d \
  -r trustRoot \
  -k ~/Library/Keychains/login.keychain-db \
  "${STAGE:?}/secrets/pki/root/pki/ca.crt"
