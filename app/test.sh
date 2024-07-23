#!/bin/bash
set -e -E
set -o pipefail

# to clean before & after
ports=({8080..8082} {8090..8091})

function nextcolor() {
  echo $(($(jobs -p | wc -l) % 6 + 31))
}

function prefix() {
  COLOR="${1:?}"
  shift 1
  PREFIX="$(echo "$@" | tr -d ';')"
  esc=$(printf '\033')
  "$@" | gsed -E "s;^(.*)$;${esc}\[${COLOR}m$PREFIX:${esc}\[0m \1;g"
}

function cleanports() {
  for port in "${ports[@]}"; do
    lsof -i ":$port" >/dev/null 2>&1 && kill -9 "$(lsof -i ":$port" | tail -n 1 | cut -d ' ' -f 2)" >/dev/null 2>&1
  done
  return 0
}

function pings() {
  for port in "${ports[@]}"; do
    until curl --fail "https://localhost:${port}/ping" >/dev/null 2>&1; do sleep 1; done
  done
}

function service() {
  SERVICENAME="${1:?}"
  if test -d "cmd/${SERVICENAME}/database"; then
    /usr/local/go/bin/go test -timeout 10s -run '^TestMigration$' "logbook/cmd/${SERVICENAME}/database" -v -count=1
  fi
  unbuffer go run "logbook/cmd/${SERVICENAME}" \
    -e local \
    -api api.yml \
    -deployment ../platform/local/deployment.yml \
    -service "cmd/${SERVICENAME}/local.yml" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function registry() {
  unbuffer go run "logbook/cmd/registry" \
    -e local \
    -api api.yml \
    -deployment ../platform/local/deployment.yml \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

function gateway() {
  GATEWAYNAME="${1:?}"
  unbuffer go run "logbook/cmd/$GATEWAYNAME" \
    -e local \
    -api api.yml \
    -deployment ../platform/local/deployment.yml \
    -discovery "../platform/local/discovery.yml" \
    -cert "../platform/local/tls/localhost.crt" \
    -key "../platform/local/tls/localhost.key"
}

trap cleanports EXIT
cleanports

prefix "$(nextcolor)" registry &
prefix "$(nextcolor)" service account &
prefix "$(nextcolor)" service objectives &
prefix "$(nextcolor)" gateway api &
prefix "$(nextcolor)" gateway internal &

pings
set +e

prefix "$(nextcolor)" unbuffer httpyac tests/api/discovery.http --all --insecure
prefix "$(nextcolor)" unbuffer httpyac tests/api/accounts.http --all --insecure
prefix "$(nextcolor)" unbuffer httpyac tests/api/objectives.http --all --insecure
