#!/usr/local/bin/bash

: "${ANSIBLE_SUDO_USER_PASSWD_HASH:?}"
: "${DIGITALOCEAN_ACCESS_TOKEN:?}"
: "${TF_VAR_DIGITALOCEAN_TOKEN:?}"
: "${TF_VAR_OVPN_AUTH_HASH:?}"
: "${TF_VAR_OVPN_AUTH_TOTP:?}"
: "${TF_VAR_OVPN_AUTH_USERNAME:?}"

(ssh-add -l | grep "do" >/dev/null) ||
  (echo "calling ssh-agent" && ssh-agent && ssh-add secrets/ssh/do)

export STAGE="${WORKSPACE:?}/platform/stage"

# shellcheck disable=SC2155
export TF_VAR_DO_SSH_FINGERPRINT="$(
  doctl compute ssh-key list --no-header --output json |
    jq -r '.[] | select(.name == "logbook") | .fingerprint' |
    tail -n 1
)"

alias ssh="ssh -F '${STAGE:?}/artifacts/ssh.conf'"
alias scp="scp -F '${STAGE:?}/artifacts/ssh.conf'"

_ssh_completion() {
  local cur opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  # prev="${COMP_WORDS[COMP_CWORD - 1]}"
  opts=$(grep '^Host' "${STAGE:?}/artifacts/ssh.conf" 2>/dev/null | grep -v '[?*]' | cut -d ' ' -f 2-)
  COMPREPLY=($(compgen -W "$opts" -- "${cur}"))
  return 0
}
complete -F _ssh_completion ssh

# MARK: Utilities

PING_URL="stage.logbook.balaasad.com:8080/api/v1.0.0/ping"

