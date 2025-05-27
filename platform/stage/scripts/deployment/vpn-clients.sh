#!/usr/bin/env bash

set -xe

: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Upload Ovpn-Auth database
# ---------------------------------------------------------------------------- #

test -f "$STAGE/provision/vpn/terraform.tfstate" && (
  # digitalocean
  jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provision/vpn/terraform.tfstate" |
    while read -r DROPLET; do
      IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
      REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"

      cd "$STAGE/secrets/pki/vpn"
      SERVERNAME="$REGION.do.vpn.logbook"
      easyrsa --batch build-server-full "$SERVERNAME" nopass

      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/ca.crt" "$VPS_SUDO_USER@$IP:ca.crt"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/issued/$SERVERNAME.crt" "$VPS_SUDO_USER@$IP:server.crt"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/private/$SERVERNAME.key" "$VPS_SUDO_USER@$IP:server.key"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn-users/crl.pem" "$VPS_SUDO_USER@$IP:crl.pem"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secrets/ovpn-auth/ovpn_auth_database.yml" "$VPS_SUDO_USER@$IP:ovpn_auth_database.yml"

      ssh -i "$STAGE/secrets/ssh/do" -n "$VPS_SUDO_USER@$IP" "sudo bash -c '
        set -xe
        
        cd /home/$VPS_SUDO_USER
        
        mv ca.crt /etc/openvpn/
        mv crl.pem /etc/openvpn/
        mv server.crt /etc/openvpn/
        mv server.key /etc/openvpn/
        
        mv ovpn_auth_database.yml /etc/openvpn/
        chmod 744 /etc/openvpn/ovpn_auth_database.yml
        chown root:root /etc/openvpn/ovpn_auth_database.yml
      '"
    done
)
