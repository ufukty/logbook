#!/usr/bin/env bash

set -xe

test -d secrets/ovpn-auth ||
  mkdir -p secrets/ovpn-auth

cd secrets/ovpn-auth

test -f ovpn_auth_database.yml &&
  rm -v ovpn_auth_database.yml

ovpn-auth register --database ovpn_auth_database.yml
