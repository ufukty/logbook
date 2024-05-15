#!/usr/local/bin/bash

_check_ssh() {
    SSH_KEY_NAME="mbp-ed"
    ssh-add -l | grep ${SSH_KEY_NAME} >/dev/null
}

_enable_ssh() {
    note "Calling ssh-agent" && ssh-agent && ssh-add
}

_check_ssh || _enable_ssh

test "$VIRTUAL_ENV" || . "$HOME/venv/bin/activate"

is_newer_than_all() {
    local TARGET="$1"
    local -a SOURCES=("${@:2}")
    test -e "$TARGET" || return 1
    for SOURCE in "${SOURCES[@]}"; do
        test -e "$SOURCE" || continue
        test "$TARGET" -ot "$SOURCE" && return 1
    done
    return 0
}
export -f is_newer_than_all

is_up_to_date() {
    is_newer_than_all "$1" $(find . -type f | grep -v "$1")
}
export -f is_up_to_date
