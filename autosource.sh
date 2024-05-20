#!/usr/local/bin/bash

export WORKSPACE
WORKSPACE="$(pwd -P)"

# MARK: Utilities

function note() { echo -e "\033[30m\033[43m\033[1m $* \033[0m"; }
function error() { echo -e "\033[38m\033[41m\033[1m $* \033[0m" >&2 && return 1; }
function red() { echo -e "\e[31m$(cat)\e[0m"; }
function green() { echo -e "\e[32m$(cat)\e[0m"; }
function yellow() { echo -e "\e[33m$(cat)\e[0m"; }
function blue() { echo -e "\e[34m$(cat)\e[0m"; }
function magenta() { echo -e "\e[35m$(cat)\e[0m"; }
function cyan() { echo -e "\e[36m$(cat)\e[0m"; }

export -f note
export -f error
export -f red
export -f green
export -f yellow
export -f blue
export -f magenta
export -f cyan

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

export -f confirm

    fi
}
export -f with-echo

note() {
    echo -e "\033[30m\033[43m\033[1m $* \033[0m"
}
export -f note

error() {
    echo -e "\033[38m\033[41m\033[1m $* \033[0m"
}
export -f error
