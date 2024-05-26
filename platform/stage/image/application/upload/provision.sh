#!/bin/bash

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${SUDO_USER:?}"
: "${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}"

# A Linux and Postgres user will be created with this name, in addition to a Postgres Database
: "${POSTGRES_USER:?}"
: "${POSTGRES_SERVER_PRIVATE_IP:?}"

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
# Function definitions
# ---------------------------------------------------------------------------- #

function app-db-tunnel() {
    sed --in-place \
        -e "s/{{SUDO_USER}}/$SUDO_USER/g" \
        -e "s/{{POSTGRES_USER}}/$POSTGRES_USER/g" \
        -e "s/{{POSTGRES_SERVER_PRIVATE_IP}}/$POSTGRES_SERVER_PRIVATE_IP/g" \
        "/etc/systemd/system/app-db-tunnel.service"

    systemctl daemon-reload
    systemctl enable app-db-tunnel
    systemctl start app-db-tunnel
}

function app-service() {
    sed --in-place \
        -e "s/<<SUDO_USER>>/${SUDO_USER:?}/g" \
        -e "s/<<POSTGRES_USER>>/${POSTGRES_USER:?}/g" \
        -e "s/<<POSTGRES_SERVER_PRIVATE_IP>>/${POSTGRES_SERVER_PRIVATE_IP:?}/g" \
        "/etc/systemd/system/app-db-tunnel.service"

    systemctl daemon-reload

    systemctl enable picarus-sync-postgres-tunnel
    systemctl start picarus-sync-postgres-tunnel

    systemctl enable picarus-sync-backend
    systemctl start picarus-sync-backend
}

function configure-iptables() {
    sed --in-place \
        -e "s/{{PRIVATE_ETHERNET_INTERFACE}}/$IPTABLES_PRIVATE_ETHERNET_INTERFACE/g" \
        "/etc/iptables/custom-rules.v4"

    systemctl restart custom-rules
}

function configure-ssh() {
    # "add public key of Postgres server to .ssh/known_hosts"
    ssh-keyscan $POSTGRES_SERVER_PRIVATE_IP >>"/home/$SUDO_USER/.ssh/known_hosts"
    ssh-keyscan $POSTGRES_SERVER_PRIVATE_IP >>"/root/.ssh/known_hosts"

    # "update .ssh directory with correct ownership and permissions"
    chmod -R 700 "/home/$SUDO_USER/.ssh"
    chown -R $SUDO_USER:$SUDO_USER "/home/$SUDO_USER/.ssh"
}

app-db-tunnel
configure-iptables
configure-ssh
