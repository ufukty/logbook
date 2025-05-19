#!/usr/bin/env bash

PS4="\n> "
set -xTve

mkdir -p "${STAGE:?}/secrets"
cd "${STAGE:?}/secrets"

mkdir -p "${STAGE:?}/secrets/ssh"
test -f "${STAGE:?}/secrets/ssh/application-server" && rm -rfv "${STAGE:?}/secrets/ssh/application-server"
ssh-keygen -a 1000 -b 4096 -C "application-server" -o -t rsa -f "${STAGE:?}/secrets/ssh/application-server" -N ''
cp "${STAGE:?}/secrets/ssh/application-server" "${STAGE:?}/image/application/upload/map/home/.ssh/application-server"
cp "${STAGE:?}/secrets/ssh/application-server.pub" "${STAGE:?}/image/database/upload/map/home/.ssh/authorized_keys"

rotate-server-cert() {
  COMMON_NAME="${1:?}"
  # https://github.com/OpenVPN/easy-rsa/blob/master/doc/EasyRSA-Renew-and-Revoke.md
  if test -f "${PKI:?}/issued/${COMMON_NAME:?}.crt"; then
    if test -f "${PKI:?}/expired/${COMMON_NAME:?}.crt"; then
      easyrsa --batch revoke-expired "${COMMON_NAME:?}" unspecified
    fi
    easyrsa --batch expire "${COMMON_NAME:?}"
    easyrsa --batch sign-req server "${COMMON_NAME:?}"
  else
    easyrsa --subject-alt-name="DNS:${COMMON_NAME:?}" --batch build-server-full "${COMMON_NAME:?}" nopass
  fi
}
PKI="${STAGE:?}/secrets/pki"

rotate-server-cert "stage.logbook.balaasad.com"

cp "${PKI:?}/issued/stage.logbook.balaasad.com.crt" \
  "${STAGE:?}/image/gateway/upload/map/etc/ssl/certs/stage.logbook.balaasad.com.crt"
cp "${PKI:?}/private/stage.logbook.balaasad.com.key" \
  "${STAGE:?}/image/gateway/upload/map/etc/ssl/private/stage.logbook.balaasad.com.key"
