#!/bin/bash
set -x
set -m

function cleanup() {
    set +x
    jobs -p | xargs -I {} kill -INT -{} || echo
}

trap cleanup EXIT

for s in cmd/*; do
    /usr/local/go/bin/go test -timeout 10s -run '^TestMigration$' "logbook/$s/database" -v -count=1
    go run "logbook/$s" \
        -a api.yml \
        -d ../platform/local.yml \
        -s "$s/testing.yml" \
        -cert "../platform/local/tls/localhost.crt" \
        -key "../platform/local/tls/localhost.key" &
done

sleep 4
httpyac examples.rest --all --insecure
