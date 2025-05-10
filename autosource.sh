#!/usr/local/bin/bash

# MARK: Compile

version() {
  echo "$(date -u +%y%m%d-%H%M%S)-$(git describe --tag --always --dirty)"
}

# MARK: Re-Deployment (only binaries for one server kind)

redeploy() {
  PROGRAM_NAME="$1" && shift
  (cd platform/stage && ./commands deploy "$PROGRAM_NAME" "$@")
}

build-redeploy() {
  PROGRAM_NAME="$1" && shift
  build "$PROGRAM_NAME"
  redeploy "$PROGRAM_NAME"
}

# MARK: API

api-summary() {
  cat be/api.http | grep HTTP/1.1 |
    cut -d ' ' -f 1-2 | awk '{ print $2, $1 }' |
    sort | awk '{ print $2, "\t", $1 }' |
    sed -E 's/(.*){{api}}(.*)/\1 \2/'
}

api-update() {
  API_GATEWAY_IP_ADDRESS="$(cat platform/stage/artifacts/deployment/service_discovery.json | jq -r '.digitalocean.fra1.services["api-gateway"][0].ipv4_address')"
  gsed --in-place "s;^@api.*;@api = http://${API_GATEWAY_IP_ADDRESS}:8080/api/v1.0.0;" be/api.http
}

ip-of() {
  PROGRAM_NAME="$1"
  cat platform/stage/deployment/service_discovery.json | jq -r ".[\"$PROGRAM_NAME\"].digitalocean[0].ipv4_address"
}
