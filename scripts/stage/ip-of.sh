#!/usr/bin/env bash

# shellcheck disable=SC2002
cat "platform/stage/deployment/service_discovery.json" |
  jq -r ".[\"$PROGRAM_NAME\"].digitalocean[0].ipv4_address"
