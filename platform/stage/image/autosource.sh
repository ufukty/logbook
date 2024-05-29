#!/usr/local/bin/bash

build-base() (
    set -x -v -e -E
    cd "${STAGE:?}/image/base"
    bash build.sh
)

build-internal() (
    set -x -v -e -E
    test "$2" == "-d" && build-base "-d"
    cd "${STAGE:?}/image/internal"
    bash build.sh
)

build-vpn() (
    set -x -v -e -E
    test "$2" == "-d" && build-base "-d"
    cd "${STAGE:?}/image/vpn"
    bash build.sh
)

build-application() (
    set -x -v -e -E
    test "$2" == "-d" && build-internal "-d"
    cd "${STAGE:?}/image/application"
    bash build.sh
)

# build-database() (
#     set -x -v -e -E
#     test "$2" == "-d" && build-internal "-d"
#     cd "${STAGE:?}/image/database"
#     bash build.sh
# )

build-gateway() (
    set -x -v -e -E
    test "$2" == "-d" && build-internal "-d"
    cd "${STAGE:?}/image/gateway"
    bash build.sh
)

build-all() (
    set -x -v -e -E
    build-base
    build-internal
    build-vpn ""
    build-application ""
    # build-database ""
    build-gateway ""
)
