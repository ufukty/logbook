#!/usr/local/bin/bash

function _autosource() {
    local PARENTDIR
    PARENTDIR="$(dirname "$1")"
    test "$PARENTDIR" && test "$PARENTDIR" != "/" && test -d "$PARENTDIR" && _autosource "$PARENTDIR"
    cd "$1" || return
    test -f "$1/autosource.sh" && echo "+ source $PWD/autosource.sh" && source "autosource.sh"
}
_autosource "$PWD"

cde() {
    set +x
    local START_DIR
    START_DIR="$(pwd -P)"
    cd "$@" || return
    set +E
    set +e
    _autosource "$PWD"
    OLDPWD="$START_DIR"
}
