#!/usr/bin/env bash

export WORKSPACE="$PWD"

# shellcheck disable=1091
test "$VIRTUAL_ENV" || . "$WORKSPACE/.venv/bin/activate"

# shellcheck disable=SC1090
. ./*.env
