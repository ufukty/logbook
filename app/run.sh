#!/bin/bash
set -e -E

function service() {
  SERVICENAME="${1:?}"
  /usr/local/go/bin/go test -timeout 10s -run '^TestMigration$' "logbook/cmd/${SERVICENAME}/database" -v -count=1
  go run "logbook/cmd/${SERVICENAME}" \
    -e local \
    -api api.yml \
    -deployment ../platform/local.yml \
    -service "cmd/${SERVICENAME}/local.yml" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function gateway() {
  go run "logbook/cmd/gateway" \
    -e local \
    -api api.yml \
    -deployment ../platform/local.yml \
    -discovery "internal/web/discovery/models/local/service_discovery.yml" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function nextcolor() {
  echo $(($(jobs -p | wc -l) % 6 + 31))
}

function prefix() {
  PREFIX="${1:?}"
  COLOR="${2:?}"
  esc=$(printf '\033')
  gsed -E "s/^(.*)$/${esc}\[${COLOR}m${PREFIX}:${esc}\[0m \1/g"
}

service account 2>&1 | prefix account "$(nextcolor)" &
service objectives 2>&1 | prefix objectives "$(nextcolor)" &
gateway 2>&1 | prefix gateway "$(nextcolor)" &

wait
