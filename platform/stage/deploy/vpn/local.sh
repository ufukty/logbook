#!/usr/bin/env bash
# shellcheck disable=SC2155,SC2029

PS4='\033[34m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "$(test "$(basename "$PWD")" == stage)"
: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Providers
# ---------------------------------------------------------------------------- #

digitalocean() {
  test -f "$STAGE/provision/vpn/terraform.tfstate" &&
    jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provision/vpn/terraform.tfstate"
}

# ---------------------------------------------------------------------------- #
# Utils
# ---------------------------------------------------------------------------- #

redact() {
  perl -0777 -pe '
    s{(BEGIN CERTIFICATE).*?(END CERTIFICATE)}{$1/REDACTED FOR LOGS/$2}gs;
    s{(BEGIN PRIVATE KEY).*?(END PRIVATE KEY)}{$1/REDACTED FOR LOGS/$2}gs;
  ' /dev/stdin
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

cd "$STAGE"

digitalocean | while read -r DROPLET; do
  PS4='\033[34m$0:$LINENO\033[0m: \033[33m$SERVER_NAME ($PUBLIC_IP)\033[0m: '

  PUBLIC_IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
  REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
  SERVER_NAME="$REGION.do.vpn.logbook"

  test -f "secrets/pki/vpn/issued/$SERVER_NAME.crt" ||
    EASYRSA_PKI="secrets/pki/vpn" easyrsa --batch build-server-full "$SERVER_NAME" nopass

  set +x
  export GATEWAY_IP="127.0.0.1" # FIXME:
  export OPENVPN_SUBNET_ADDRESS="$(sed 's;//.*;;g' <"$STAGE/config/digitalocean.jsonc" | jq --arg region "$REGION" -r '.vpn[$region]')"
  export UNBOUND_ADDRESS="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address_private')"
  export VPC_CIDR="$(doctl vpcs get "$(echo "$DROPLET" | jq -r '.attributes.vpc_uuid')" --output json | jq -r '.[0].ip_range')"
  export VPC_RANGE_ADDRESS="$(echo "$VPC_CIDR" | perl -nE 'say $1 if /^(.*)\//')"
  export VPC_RANGE_MASK="$(echo "$VPC_CIDR" | perl -nE 'say $1 if /\/(\d{2})$/')"
  export ROOT_CA_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/root/ca.crt")"
  export VPN_CA_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/vpn/ca.crt")"
  export VPN_SERVER_CERT="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/vpn/issued/$SERVER_NAME.crt")"
  export VPN_SERVER_KEY="$(awk '/BEGIN/,/END/' "$STAGE/secrets/pki/vpn/private/$SERVER_NAME.key")"
  set -x

  SSH_ARGS=(-i "$STAGE/secrets/ssh/do" -o ControlMaster=auto -o ControlPath="$(mktemp -u)" -o ControlPersist=300)

  TEMPLATES="$STAGE/deploy/vpn/template"
  find "$TEMPLATES" -type f | while read -r TEMPLATE; do
    perl -nE 'say $1 if /\${(.*?)}/' <"$TEMPLATE" | sort | uniq | while read -r ENVVAR; do printenv "$ENVVAR" >/dev/null; done
    (colordiff <(cat "$TEMPLATE") <(envsubst <"$TEMPLATE") || true) | redact
    envsubst <"$TEMPLATE" | ssh "${SSH_ARGS[@]}" "$VPS_SUDO_USER@$PUBLIC_IP" sudo tee "${TEMPLATE//"$TEMPLATES"/}" >/dev/null
  done

  scp "${SSH_ARGS[@]}" \
    "secrets/pki/vpn-users/crl.pem" \
    "secrets/ovpn-auth/ovpn_auth_database.yml" \
    "deploy/vpn/remote.sh" \
    "$VPS_SUDO_USER@$PUBLIC_IP:"

  scp "${SSH_ARGS[@]}" "secrets/vpn/tls-crypt/do-$REGION.key" "$VPS_SUDO_USER@$PUBLIC_IP:tls-crypt.key"

  ssh "${SSH_ARGS[@]}" -n "$VPS_SUDO_USER@$PUBLIC_IP" "sudo --preserve-env bash remote.sh && rm -rfv *"
done
