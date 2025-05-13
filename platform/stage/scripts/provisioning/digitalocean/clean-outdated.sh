#!/usr/local/bin/bash

function note() { echo -e "\033[30m\033[43m\033[1m $* \033[0m"; }
function error() { echo -e "\033[38m\033[41m\033[1m $* \033[0m" >&2 && return 1; }
function red() { echo -e "\e[31m$(cat)\e[0m"; }
function green() { echo -e "\e[32m$(cat)\e[0m"; }
function yellow() { echo -e "\e[33m$(cat)\e[0m"; }
function blue() { echo -e "\e[34m$(cat)\e[0m"; }
function magenta() { echo -e "\e[35m$(cat)\e[0m"; }
function cyan() { echo -e "\e[36m$(cat)\e[0m"; }

function confirm() {
  read -p "Proceed? (y) " -n 1 USER_INPUT
  echo
  case $USER_INPUT in
  [Yy]*)
    return 0
    ;;
  esac
  return 1
}

set -E
ALL="$(doctl compute snapshot list --format Name,ID --no-header | grep build | grep -v base | sort -r | gsed -E 's/[ ]+/ /g')"
TOKEEP="$(echo "$ALL" | gsed -E 's/(build_([^_]+)_([0-9_]{17}).*)/\1 \2/g' | sort -r | uniq -f 2 | cut -d ' ' -f 1)"
OUTDATED="$(echo "$ALL" | grep -v "$TOKEEP")"
OUTDATED_NAMEs="$(echo "$OUTDATED" | cut -d ' ' -f 1)"
OUTDATED_IDs="$(echo "$OUTDATED" | cut -d ' ' -f 2)"

test -z "$OUTDATED_IDs" && echo "nothing to clean" && return 0

echo "keep:" | green
echo "$TOKEEP" | green
echo "delete:" | red
echo "$OUTDATED_NAMEs" | red

confirm "$ALL" "$OUTDATED_IDs" && echo "$OUTDATED_IDs" | xargs --verbose -n 1 doctl compute image delete -f
