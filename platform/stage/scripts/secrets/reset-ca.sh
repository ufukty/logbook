#!/usr/bin/env bash

PS4="\n> "
set -xTveE

mkdir -p "${STAGE:?}/secrets"
cd "${STAGE:?}/secrets"
test -d pki && rm -rfv pki

easyrsa init-pki soft
easyrsa --batch --req-cn="Logbook Stage Environment CA" build-ca nopass

security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db "${STAGE:?}/secrets/pki/ca.crt" # macos keychain
