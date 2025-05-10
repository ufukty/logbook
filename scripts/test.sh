#!/usr/local/bin/bash

set -xeEo pipefail
PS4='\033[31m$0:$LINENO: \033[0m'

test -d .git || (echo "Run from root folder" && exit 1)

for port in {8080..8082}; do
  until curl --fail "https://localhost:${port}/ping" >/dev/null 2>&1; do sleep 1; done
done

set +e

echo "starting tests..." | prefix "$YELLOW" test
unbuffer httpyac tests/api/discovery.http --all --insecure 2>&1 | prefix "$CYAN" discovery
unbuffer httpyac tests/api/accounts.http --all --insecure 2>&1 | prefix "$BLUE" accounts
unbuffer httpyac tests/api/objectives.http --all --insecure 2>&1 | prefix "$MAGENTA" objectives
