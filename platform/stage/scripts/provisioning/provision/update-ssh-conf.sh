#!/usr/bin/env bash

set -xe

: "${STAGE:?}"

cat "$STAGE"/artifacts/ssh.conf.d/* >"$STAGE/artifacts/ssh.conf"
