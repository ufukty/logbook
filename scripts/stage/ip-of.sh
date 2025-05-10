#!/usr/local/bin/bash

cat "platform/stage/deployment/service_discovery.json" |
  jq -r ".[\"$PROGRAM_NAME\"].digitalocean[0].ipv4_address"
