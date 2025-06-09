#!/usr/bin/env bash
# shellcheck disable=SC2155

PS4='\033[32m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${STAGE:?}"

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  # shellcheck disable=SC2002,SC2046
  cat $(find find "$STAGE/provision" -name 'terraform.tfstate') |
    jq -c 'select(.resources | length > 0) | .resources.[] | select(.type == "digitalocean_droplet").instances.[]'
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

mkdir -p "$STAGE/artifacts/vpn-clients"

# shellcheck disable=SC2002
digitalocean | while read -r HOST; do
  IP="$(echo "$HOST" | jq -r '.attributes.ipv4_address')"
  REGION="$(echo "$HOST" | jq -r '.attributes.region')"

  cat "$STAGE/config/vpn-users" | while read -r VPN_USER; do
    PROFILE_NAME="$VPN_USER@$REGION.do.vpn.logbook"

    test -f "$STAGE/secrets/pki/vpn-users/reqs/$PROFILE_NAME.req" ||
      EASYRSA_PKI="$STAGE/secrets/pki/vpn-users" EASYRSA_BATCH="yes" easyrsa build-client-full "$PROFILE_NAME" nopass

    set +x
    export PUBLIC_IP="$IP"
    export EASYRSA_SERVER_NAME="$REGION.do.vpn.logbook"
    export VPN_USERS_CA_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/vpn-users/ca.crt")"
    export ROOT_CA_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/root/ca.crt")"
    export VPN_USER_KEY="$(cat "$STAGE/secrets/pki/vpn-users/private/$PROFILE_NAME.key")"
    export VPN_USER_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/vpn-users/issued/$PROFILE_NAME.crt")"
    export TLS_CRYPT_KEY="$(cat "$STAGE/secrets/vpn/tls-crypt/do-$REGION.key")"
    set -x

    TEMPLATE="$STAGE/scripts/artifacts/template/vpn-client.ovpn"
    perl -nE 'say $1 if /\${(.*?)}/' <"$TEMPLATE" | sort | uniq | while read -r ENVVAR; do printenv "$ENVVAR" >/dev/null; done
    envsubst <"$TEMPLATE" >"$STAGE/artifacts/vpn-clients/$PROFILE_NAME.ovpn"
  done
done
