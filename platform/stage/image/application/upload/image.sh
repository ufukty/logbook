#!/bin/bash

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}"

# ---------------------------------------------------------------------------- #
# Commons
# ---------------------------------------------------------------------------- #

PS4="\033[32m/\$(basename \"\${BASH_SOURCE}\"):\${LINENO}\033[0m\033[33m\${FUNCNAME[0]:+/\${FUNCNAME[0]}():}\033[0m "
set -x -T -v -e -E

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
# Tasks
# ---------------------------------------------------------------------------- #

function iptables-configure() {
    sed --in-place \
        -e "s/{{PRIVATE_ETHERNET_INTERFACE}}/${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}/g" \
        "/etc/iptables/iptables-rules.v4"
    systemctl restart iptables-activation
}

function configure-logging() {
    chown -R syslog:root "/var/log/logbook-app"
    chown syslog:root "/var/log/syslog"
    chmod 640 "/var/log/syslog"
    systemctl restart rsyslog
    # TODO: "managing logrotate"
}

function ssh-dir() {
    chmod -R 700 "/home/${SUDO_USER:?}/.ssh"
    chown -R "${SUDO_USER:?}:${SUDO_USER:?}" "/home/${SUDO_USER:?}/.ssh"
}

# ---------------------------------------------------------------------------- #
# Main
# ---------------------------------------------------------------------------- #

assert_sudo
restart_journald
remove_password_change_requirement
check_tun_availability
wait_cloud_init
apt_update

deploy_provisioner_files

iptables-configure
configure-logging
ssh-dir
