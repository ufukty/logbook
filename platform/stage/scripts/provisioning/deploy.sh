#!/usr/local/bin/bash

export PROGRAM_NAME="$1" && shift

cd "${STAGE:?}/deployment"
if test -z "$PROGRAM_NAME"; then
  ansible-playbook --forks="20" playbook.yml
else
  ansible-playbook --forks="20" --limit="$PROGRAM_NAME" --tags="redeploy" playbook.yml
fi

test "$(curl -sSL "${PING_URL:?}")" = "pong" || error "API gateway didn't pong to ping"

