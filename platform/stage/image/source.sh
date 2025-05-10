#!/usr/local/bin/bash

build-base() (
    set -x -v -e -E
    cd "${STAGE:?}/image/base"
    bash build.sh "$@"
)

build-internal() (
    set -x -v -e -E
    test "$2" == "-d" && shift && build-base "-d"
    cd "${STAGE:?}/image/internal"
    bash build.sh "$@"
)

build-vpn() (
    set -x -v -e -E
    test "$2" == "-d" && shift && build-base "-d"
    cd "${STAGE:?}/image/vpn"
    bash build.sh "$@"
)

build-application() (
    set -x -v -e -E
    test "$2" == "-d" && shift && build-internal "-d"
    cd "${STAGE:?}/image/application"
    bash build.sh "$@"
)

# build-database() (
#     set -x -v -e -E
#     test "$2" == "-d" && shift && build-internal "-d"
#     cd "${STAGE:?}/image/database"
#     bash build.sh "$@"
# )

build-gateway() (
    set -x -v -e -E
    test "$2" == "-d" && shift && build-internal "-d"
    cd "${STAGE:?}/image/gateway"
    bash build.sh "$@"
)

# build-all [ -B ]
build-all() (
    PS4='\033[31m$(basename "${BASH_SOURCE}"):${LINENO}\033[0m\033[33m${FUNCNAME[0]:+/${FUNCNAME[0]}():}\033[0m '
    set -x -v -e -E
    build-base "$@"
    build-internal "$@"
    build-vpn "" "$@"
    build-application "$@"
    # build-database "$@"
    build-gateway "$@"
)
