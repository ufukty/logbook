#!/usr/local/bin/bash

DISCOVERY="${STAGE:?}/artifacts/deployment/service_discovery.json"

touch "$DISCOVERY"
ADDRESSES="$(jq -r '.digitalocean.fra1.services[] | .[] | .ipv4_address_private' <"$DISCOVERY")"
test -z "$ADDRESSES" && return
echo "$ADDRESSES" | while read ADDRESS; do
  ssh-keygen -R "$ADDRESS" >/dev/null 2>&1
  ssh-keyscan "$ADDRESS" >>~/.ssh/known_hosts 2>/dev/null
done
