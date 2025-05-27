#!/usr/bin/env bash

PS4='\033[32m$0:$LINENO\033[0m: '
set -xe

: "${VPS_SUDO_USER:?}"

test -f "$STAGE/provision/vpn/terraform.tfstate" && (

  # ---------------------------------------------------------------------------- #
  # DO regions
  # ---------------------------------------------------------------------------- #

  jq -c '.resources.[] | select(.type == "digitalocean_droplet").instances.[]' <"$STAGE/provision/vpn/terraform.tfstate" |
    while read -r DROPLET; do
      IP="$(echo "$DROPLET" | jq -r '.attributes.ipv4_address')"
      REGION="$(echo "$DROPLET" | jq -r '.attributes.region')"
      SERVERNAME="$REGION.do.vpn.logbook"

      EASYRSA_PKI="$STAGE/secrets/pki/vpn" easyrsa --batch build-server-full "$SERVERNAME" nopass

      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/ca.crt" "$VPS_SUDO_USER@$IP:ca.crt"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/issued/$SERVERNAME.crt" "$VPS_SUDO_USER@$IP:server.crt"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn/private/$SERVERNAME.key" "$VPS_SUDO_USER@$IP:server.key"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secret/pki/vpn-users/crl.pem" "$VPS_SUDO_USER@$IP:crl.pem"
      scp -i "$STAGE/secrets/ssh/do" "$STAGE/secrets/ovpn-auth/ovpn_auth_database.yml" "$VPS_SUDO_USER@$IP:ovpn_auth_database.yml"

      ssh -i "$STAGE/secrets/ssh/do" -n "$VPS_SUDO_USER@$IP" "sudo bash -c '
        set -xe
        
        cd /home/$VPS_SUDO_USER
        mv ca.crt crl.pem server.crt server.key ovpn_auth_database.yml /etc/openvpn/

        cd /etc/openvpn
        chmod 644 crl.pem
        chmod 744 ovpn_auth_database.yml
        chown root:root ovpn_auth_database.yml
      '"
    done
)
