#!/usr/local/bin/bash

function autosource() {
    local TARGET
    TARGET="${1:-"$PWD"}"
    local PARENTDIR
    PARENTDIR="$(dirname "$TARGET")"
    test "$PARENTDIR" && test "$PARENTDIR" != "/" && test -d "$PARENTDIR" && autosource "$PARENTDIR"
    cd "$TARGET" || return
    test -f "$TARGET/autosource.sh" && echo "+ source $PWD/autosource.sh" && source "autosource.sh"
}
