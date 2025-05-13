#!/usr/local/bin/bash

cat "${STAGE:?}/artifacts/ssh.conf.d/"* >"${STAGE:?}/artifacts/ssh.conf"
