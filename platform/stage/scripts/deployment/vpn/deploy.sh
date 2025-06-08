#!/usr/bin/env bash
# shellcheck disable=SC2155,SC2029

PS4='\033[34m$0:$LINENO:\033[0m '
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
  test -f "$STAGE/provisioning/vpn/terraform.tfstate" &&
    jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provisioning/vpn/terraform.tfstate"
}

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

cd "$STAGE"

digitalocean | while read -r DROPLET; do
  PS4='\033[34m$0:$LINENO: \033[33m$PUBLIC_IP:\033[0m '

  export PUBLIC_IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
  export PRIVATE_IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address_private')"
  export REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
  export SERVER_NAME="$REGION.do.vpn.logbook"
  export OPENVPN_SUBNET_ADDRESS="$(sed 's;//.*;;g' <"$STAGE/config/digitalocean.jsonc" | jq --arg region "$REGION" -r '.vpn[$region]')"

  test -f "secrets/pki/vpn/issued/$SERVER_NAME.crt" ||
    EASYRSA_PKI="secrets/pki/vpn" easyrsa --batch build-server-full "$SERVER_NAME" nopass

  SSH_ARGS=(-i "$STAGE/secrets/ssh/do" -o ControlMaster=auto -o ControlPath="$(mktemp -u)" -o ControlPersist=300)

  TEMPLATES="$STAGE/scripts/deployment/vpn/template"
  find "$TEMPLATES" -type f | while read -r TEMPLATE; do
    envsubst <"$TEMPLATE" | ssh "${SSH_ARGS[@]}" "$VPS_SUDO_USER@$PUBLIC_IP" sudo tee "${TEMPLATE//"$TEMPLATES"/}" >/dev/null
  done

  scp "${SSH_ARGS[@]}" \
    "secrets/pki/vpn/ca.crt" \
    "secrets/pki/vpn-users/crl.pem" \
    "secrets/ovpn-auth/ovpn_auth_database.yml" \
    "scripts/deployment/vpn/remote.sh" \
    "$VPS_SUDO_USER@$PUBLIC_IP:"

  scp "${SSH_ARGS[@]}" "secrets/pki/vpn/issued/$SERVER_NAME.crt" "$VPS_SUDO_USER@$PUBLIC_IP:server.crt"
  scp "${SSH_ARGS[@]}" "secrets/pki/vpn/private/$SERVER_NAME.key" "$VPS_SUDO_USER@$PUBLIC_IP:server.key"

  ssh "${SSH_ARGS[@]}" -n "$VPS_SUDO_USER@$PUBLIC_IP" "sudo bash remote.sh && rm -rfv *"
done
