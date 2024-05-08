#!/bin/bash

# ---------------------------------------------------------------------------- #
# Commons
# ---------------------------------------------------------------------------- #

PS4="\033[32m/\$(basename \"\${BASH_SOURCE}\"):\${LINENO}\033[0m\033[33m\${FUNCNAME[0]:+/\${FUNCNAME[0]}():}\033[0m "
set -x
set -v
set -e
set -E

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
# Definitions
# ---------------------------------------------------------------------------- #-

function install_utilities() {
    retry apt-get install -y ipcalc
}

function install_openvpn() {
    retry apt-get install -y ca-certificates gnupg openvpn iptables openssl wget ca-certificates curl
    test -d /etc/openvpn/easy-rsa && rm -rf /etc/openvpn/easy-rsa/*
    return "0"
}

function install_easy_rsa() (
    cd "$PROVISIONER_FILES/dependencies"
    mkdir -p /etc/openvpn/easy-rsa
    tar xzf EasyRSA-3.1.7.tgz --strip-components=1 --directory /etc/openvpn/easy-rsa
)

function ovpn_auth() {
    chmod 755 /etc/openvpn/ovpn-auth-v1.0.4-linux-x64
    chown root:root /etc/openvpn/ovpn-auth-v1.0.4-linux-x64

    chmod 744 /etc/openvpn/ovpn_auth_database.yml
    chown root:root /etc/openvpn/ovpn_auth_database.yml
}

function install_unbound() {
    retry apt-get install -y unbound
}

# ---------------------------------------------------------------------------- #
# Main
# ---------------------------------------------------------------------------- #

assert_sudo
restart_journald
check_tun_availability
wait_cloud_init
apt_update

install_utilities
install_openvpn
install_easy_rsa
install_unbound

deploy_provisioner_files
ovpn_auth
