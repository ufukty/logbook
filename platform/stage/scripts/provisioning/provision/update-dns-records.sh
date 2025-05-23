#!/usr/bin/env bash

set -xe

GATEWAY_IP="$(
  jq -r '.digitalocean.fra1.services["gateway"][0].ipv4_address_private' <"$STAGE/artifacts/deployment/service_discovery.json"
)"
test -z "$GATEWAY_IP" && return

# shellcheck disable=SC2087
ssh fra1-vpn -i "$STAGE/secrets/ssh/do" sudo bash -s <<HERE
set -xe
sed "s/{{GATEWAY_IP}}/$GATEWAY_IP/g" /etc/unbound/unbound.conf.tmpl.d/custom.conf >/etc/unbound/unbound.conf.d/custom.conf
systemctl restart unbound
HERE

sudo killall mDNSResponder{,Helper}
