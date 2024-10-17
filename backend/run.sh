#!/bin/bash
set -e -E
set -o pipefail

function cleanports() {
  # to clean before & after
  echo "cleanports is running..."
  ports=({8080..8099})
  for port in "${ports[@]}"; do
    lsof -i ":$port" >/dev/null 2>&1 && kill -9 "$(lsof -i ":$port" | tail -n 1 | cut -d ' ' -f 2)" >/dev/null 2>&1
  done
  return 0
}

function service() {
  SERVICENAME="${1:?}"
  if test -d "cmd/${SERVICENAME}/database"; then
    /usr/local/go/bin/go test -timeout 10s -run '^TestMigration$' "logbook/cmd/${SERVICENAME}/database" -v -count=1
  fi
  go run "logbook/cmd/${SERVICENAME}" \
    -ip localhost \
    -api api.yml \
    -service "cmd/${SERVICENAME}/local.yml" \
    -deployment "../platform/local/deployment.yml" \
    -internal "../platform/local/registryfile.internalgateway.json" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function registry() {
  go run "logbook/cmd/registry" \
    -api api.yml \
    -deployment "../platform/local/deployment.yml" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function api-gateway() {
  go run "logbook/cmd/api" \
    -api api.yml \
    -deployment "../platform/local/deployment.yml" \
    -internal "../platform/local/registryfile.internalgateway.json" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function internal-gateway() {
  go run "logbook/cmd/internal" \
    -api api.yml \
    -deployment "../platform/local/deployment.yml" \
    -registry "../platform/local/registryfile.registryservice.json" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

trap cleanports EXIT
cleanports

registry 2>&1 | prefix "$BLUE" registry &
internal-gateway 2>&1 | prefix "$CYAN" internal &
service account 2>&1 | prefix "$RED" account &
service objectives 2>&1 | prefix "$YELLOW" objectives &
api-gateway 2>&1 | prefix "$MAGENTA" api-gateway &

wait
