#!/usr/bin/env bash

set -xe

discovery="$STAGE/artifacts/deployment/service_discovery.json"
test -f "$discovery" || exit

addresses="$(
  jq -r '.digitalocean.fra1.services[] | .[] | .ipv4_address_private' <"$discovery"
)"
test -z "$addresses" && exit

echo "$addresses" | while read -r ADDRESS; do
  ssh-keygen -R "$ADDRESS" >/dev/null 2>&1
  ssh-keyscan "$ADDRESS" >>~/.ssh/known_hosts 2>/dev/null
done
