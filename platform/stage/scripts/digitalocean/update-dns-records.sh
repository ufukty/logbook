#!/usr/local/bin/bash

set -xeo pipefail

GATEWAY_IP="$(
  cat "${STAGE:?}/artifacts/deployment/service_discovery.json" |
    jq -r '.digitalocean.fra1.services["gateway"][0].ipv4_address_private'
)"
test -z "$GATEWAY_IP" && return

ssh -t fra1-vpn "sudo bash -c 'sed --in-place \"s;{{GATEWAY_IP}};${GATEWAY_IP:?};g\" /etc/unbound/unbound.conf.tmpl.d/custom.conf > /etc/unbound/unbound.conf.d/custom.conf && systemctl restart unbound'"

sudo killall mDNSResponder{,Helper}
