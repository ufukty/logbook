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
# EasyRSA Vars
# ---------------------------------------------------------------------------- #

export EASYRSA_ALGO="ec"
export EASYRSA_BATCH="yes"
export EASYRSA_CURVE="prime256v1"
export EASYRSA_DN="org" # enables REQ

export EASYRSA_REQ_CITY="NA"
export EASYRSA_REQ_COUNTRY="TR"
export EASYRSA_REQ_EMAIL=""
export EASYRSA_REQ_ORG="Logbook"
export EASYRSA_REQ_OU="Platform Engineering"
export EASYRSA_REQ_PROVINCE=""

# ---------------------------------------------------------------------------- #
# Root CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa --req-cn="Logbook Stage Root CA" build-ca nopass

# ---------------------------------------------------------------------------- #
# Web Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/web" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/web" easyrsa --req-cn="Logbook Stage Web CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa import-req "$STAGE/secrets/pki/web/reqs/ca.req" web
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa sign-req ca web

cp "$STAGE/secrets/pki/root/issued/web.crt" \
  "$STAGE/secrets/pki/web/ca.crt"

# ---------------------------------------------------------------------------- #
# Vpn Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa --req-cn="Logbook Stage VPN CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa import-req "$STAGE/secrets/pki/vpn/reqs/ca.req" vpn
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa sign-req ca vpn

cp "$STAGE/secrets/pki/root/issued/vpn.crt" \
  "$STAGE/secrets/pki/vpn/ca.crt"

# ---------------------------------------------------------------------------- #
# Vpn Users Intermediate CA
# ---------------------------------------------------------------------------- #

EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" easyrsa init-pki
EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" easyrsa --req-cn="Logbook Stage VPN Users CA" build-ca nopass subca

EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa import-req "$STAGE/secrets/pki/vpn-users/reqs/ca.req" vpn-users
EASYRSA_PKI="$STAGE/secrets/pki/root" easyrsa sign-req ca vpn-users

cp "$STAGE/secrets/pki/root/issued/vpn-users.crt" \
  "$STAGE/secrets/pki/vpn-users/ca.crt"

EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" EASYRSA_CRL_DAYS=3650 easyrsa gen-crl

# ---------------------------------------------------------------------------- #
# Debug
# ---------------------------------------------------------------------------- #

openssl x509 -in "$STAGE/secrets/pki/root/ca.crt" -text -noout
openssl x509 -in "$STAGE/secrets/pki/web/ca.crt" -text -noout
openssl x509 -in "$STAGE/secrets/pki/vpn/ca.crt" -text -noout
openssl x509 -in "$STAGE/secrets/pki/vpn-users/ca.crt" -text -noout
