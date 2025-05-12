#!/usr/local/bin/bash

(ssh-add -l | grep logbook >/dev/null) ||
  (echo "calling ssh-agent" && ssh-agent && ssh-add)

# shellcheck disable=1091
test "$VIRTUAL_ENV" || . "$WORKSPACE/.venv/bin/activate"

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
  # shellcheck disable=2046
  is_newer_than_all "$1" $(find . -type f | grep -v "$1")
}
export -f is_up_to_date
