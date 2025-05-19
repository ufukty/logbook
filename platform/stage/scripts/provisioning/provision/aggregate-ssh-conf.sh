#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

cat "${STAGE:?}/artifacts/ssh.conf.d/"* >"${STAGE:?}/artifacts/ssh.conf"
