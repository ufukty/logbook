#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

PREFIX="${1:?Pass image prefix as the first argument.}"

ALL="$(doctl compute snapshot list --output json | jq --arg prefix "$PREFIX" '[.[] | select(.name | startswith($prefix))]')"
LATEST="$(echo "$ALL" | jq -r 'max_by(.created_at).id')"
OUTDATED="$(echo "$ALL" | jq -r --arg latest_id "$LATEST" '.[] | select(.id != $latest_id) | .id')"

# shellcheck disable=SC2086
doctl compute snapshot delete $OUTDATED --force

doctl compute snapshot list --output json | jq --arg prefix "$PREFIX" '.[] | select(.name | startswith($prefix))'
