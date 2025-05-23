#!/usr/bin/env bash

autosource() {
  local TARGET
  TARGET="${1:-"$PWD"}"
  local PARENTDIR
  PARENTDIR="$(dirname "$TARGET")"
  test "$PARENTDIR" && test "$PARENTDIR" != "/" && test -d "$PARENTDIR" && autosource "$PARENTDIR"
  cd "$TARGET" || return
  # shellcheck disable=SC1091
  test -f "$TARGET/source.sh" && echo "+ source $PWD/source.sh" && source "source.sh"
}
