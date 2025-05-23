#!/usr/bin/env bash

test "$(basename "$PWD")" == "stage" || (echo "Run from the stage folder." && exit 1)

set -x

# MARK: Cloud providers

ssh-keygen -a 100 -t ed25519 -C logbook -f secrets/ssh/do
# ssh-keygen -a 100 -t ed25519 -C logbook -f secrets/ssh/aws
# ssh-keygen -a 100 -t ed25519 -C logbook -f secrets/ssh/gcp

chmod 600 secrets/ssh/*

# shellcheck disable=SC2046
doctl compute ssh-key list --no-header --output json |
  jq -r '.[] | select(.name == "logbook") | .id' |
  xargs -I {} -n 1 doctl compute ssh-key delete {} --force

doctl compute ssh-key import logbook --public-key-file secrets/ssh/do.pub
