#!/usr/bin/env bash

PS4='\033[33m$0:$LINENO\033[0m '
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

mkdir -p "$STAGE/secrets/vpn/tls-crypt"

for REGION in sfo2 sfo3 tor1 nyc1 nyc3 lon1 ams3 fra1 blr1 sgp1; do
  openvpn --genkey secret "$STAGE/secrets/vpn/tls-crypt/do-$REGION.key"
done
