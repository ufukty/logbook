#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

export PROGRAM_NAME="${1:?}"

cd "${STAGE:?}/deployment"
if test -z "$PROGRAM_NAME"; then
  ansible-playbook --forks="20" playbook.yml
else
  ansible-playbook --forks="20" --limit="$PROGRAM_NAME" --tags="redeploy" playbook.yml
fi

test "$(curl -sSL "${PING_URL:?}")" = "pong" || error "API gateway didn't pong to ping"
