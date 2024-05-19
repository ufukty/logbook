#!/bin/bash

# ---------------------------------------------------------------------------- #
# Commons
# ---------------------------------------------------------------------------- #

PS4='\033[32m$(basename "${BASH_SOURCE}"):${LINENO}\033[0m\033[33m${FUNCNAME[0]:+/${FUNCNAME[0]}():}\033[0m '
set -v -x -e -E

PROVISIONER_FILES="$(pwd -P)"

function retry() {
    local COUNTER=0
    until "$@"; do
        EC=$?
        COUNTER=$((COUNTER + 1))
        test $COUNTER -ge 60 && exit $EC
        sleep 10
    done
}

function apt_update() { retry apt-get update; }
function restart_journald() { systemctl restart systemd-journald; }
function assert_sudo() { test $EUID -eq 0; }
function remove_password_change_requirement() { sed --in-place -E 's/root:(.*):0:0:(.*):/root:\1:18770:0:\2:/g' /etc/shadow; }
function wait_cloud_init() { cloud-init status --wait >/dev/null; }
function check_tun_availability() { test -e /dev/net/tun; }
function deploy_provisioner_files() {
    chmod 700 -R "$PROVISIONER_FILES/map"
    chown root:root -R "$PROVISIONER_FILES/map"
    rsync --verbose --recursive "$PROVISIONER_FILES/map/" "/"
    rm -rfv "$PROVISIONER_FILES/map"
}

export DEBIAN_FRONTEND=noninteractive

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}"
: "${APP_USER:?}"

# ---------------------------------------------------------------------------- #
# Tasks
# ---------------------------------------------------------------------------- #

function create-app-user() {
    adduser --disabled-password --gecos "" "${APP_USER:?}"
}

function iptables_configure() {
    sed --in-place \
        -e "s/{{PRIVATE_ETHERNET_INTERFACE}}/${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}/g" \
        "/etc/iptables/iptables-rules.v4"

    chmod 644 /etc/systemd/system/iptables-activation.service
    chown root:root /etc/systemd/system/iptables-activation.service

    systemctl daemon-reload
    systemctl enable iptables-activation
    systemctl start iptables-activation
    systemctl status iptables-activation
}

function fail2ban_configure() {
    systemctl restart fail2ban
}

# ---------------------------------------------------------------------------- #
# Main
# ---------------------------------------------------------------------------- #

assert_sudo
restart_journald
remove_password_change_requirement
wait_cloud_init

deploy_provisioner_files

create-app-user
iptables_configure
fail2ban_configure
