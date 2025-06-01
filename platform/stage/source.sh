#!/usr/bin/env bash
# shellcheck disable=SC2155

: "${DIGITALOCEAN_ACCESS_TOKEN:?}"
: "${TF_VAR_DIGITALOCEAN_TOKEN:?}"

(ssh-add -l | grep -q "$(ssh-keygen -lf secrets/ssh/do.pub)") ||
  (echo "calling ssh-agent" && ssh-agent && ssh-add secrets/ssh/do)

export STAGE="${WORKSPACE:?}/platform/stage"

export DO_SSH_FINGERPRINT="$(ssh-keygen -lf "$STAGE/secrets/ssh/do" -E md5 | perl -nE 'say $1 if /(?<=MD5:)([^\s]+)/')"
export DO_SSH_KEY_ID="$(doctl compute ssh-key get "$DO_SSH_FINGERPRINT" --output json | jq -r '.[0].id')"

export DO_SSH_PUBKEY="$(cat secrets/ssh/do.pub)"

export PING_URL="stage.logbook.balaasad.com:8080/api/v1.0.0/ping"

# shellcheck disable=SC1091
. .env

# env files to declare:
: "${VPS_SUDO_USER:?}"
: "${VPS_SUDO_USER_PASSWD_HASH:?}"

# terraform aliases
export TF_VAR_DO_SSH_FINGERPRINT="$DO_SSH_FINGERPRINT"
