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

# MARK: Compile

version() {
    echo "$(date -u +%y%m%d-%H%M%S)-$(git describe --tag --always --dirty)"
}

build() {
    set -E
    set -o pipefail

    PROGRAM_NAMES="$@"
    test "$PROGRAM_NAMES" || PROGRAM_NAMES="$(cd "$WORKSPACE/app/cmd" && find . -name 'main.go' | cut -d '/' -f 2 | sort | uniq)"

    VERSION="$(version)"
    mkdir -p "${WORKSPACE}/build/${VERSION}/"{darwin,linux}
    echo "${VERSION}" | blue

    set +E
    cd "${WORKSPACE}/app"
    for PROGRAM_NAME in $PROGRAM_NAMES; do
        echo "${PROGRAM_NAME}" | green
        GOOS=darwin GOARCH=amd64 go build -o "${WORKSPACE}/build/${VERSION}/darwin/$PROGRAM_NAME" "./cmd/$PROGRAM_NAME"
        GOOS=linux GOARCH=amd64 go build -o "${WORKSPACE}/build/${VERSION}/linux/$PROGRAM_NAME" "./cmd/$PROGRAM_NAME"
    done
}

lastbuild() (
    PROGRAM_NAME="$1" && shift
    cd "$WORKSPACE"
    if test -z "$PROGRAM_NAME"; then
        find build -type d -maxdepth 3 -mindepth 2 | sort | tail -n 1 | cut -f 2 -d '/'
    else
        find build -name "*$PROGRAM_NAME*" | sort | tail -n 1 | cut -f 2 -d '/'
    fi
)

lastbuildpath() {
    PROGRAM_NAME="$1" && shift
    ARCH="${1:-"darwin"}" && shift
    PROGRAM_LAST_BUILD="$(lastbuild $PROGRAM_NAME)"
    echo "app/build/$PROGRAM_LAST_BUILD/$ARCH/$PROGRAM_NAME"
}

# MARK: Re-Deployment (only binaries for one server kind)

redeploy() {
    PROGRAM_NAME="$1" && shift
    (cd platform/stage && ./commands deploy "$PROGRAM_NAME" "$@")
}

build-redeploy() {
    PROGRAM_NAME="$1" && shift
    build "$PROGRAM_NAME"
    redeploy "$PROGRAM_NAME"
}

# MARK: API

api-summary() {
    cat app/api.http | grep HTTP/1.1 |
        cut -d ' ' -f 1-2 | awk '{ print $2, $1 }' |
        sort | awk '{ print $2, "\t", $1 }' |
        sed -E 's/(.*){{api}}(.*)/\1 \2/'
}

api-update() {
    API_GATEWAY_IP_ADDRESS="$(cat platform/stage/artifacts/deployment/service_discovery.json | jq -r '.digitalocean.fra1.services["api-gateway"][0].ipv4_address')"
    gsed --in-place "s;^@api.*;@api = http://${API_GATEWAY_IP_ADDRESS}:8080/api/v1.0.0;" app/api.http
}

ip-of() {
    PROGRAM_NAME="$1"
    cat platform/stage/deployment/service_discovery.json | jq -r ".[\"$PROGRAM_NAME\"].digitalocean[0].ipv4_address"
}
