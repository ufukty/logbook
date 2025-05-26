#!/bin/bash

PS4='\033[33m''\D{%H:%M:%S} provision:${LINENO}:''\033[0m '
set -xe

# ---------------------------------------------------------------------------- #
# Wait
# ---------------------------------------------------------------------------- #

timeout 180 bash -c "until stat /var/lib/cloud/instance/boot-finished 2>/dev/null; do sleep 2; done"

# ---------------------------------------------------------------------------- #
# Pre checks
# ---------------------------------------------------------------------------- #

test $EUID -eq 0 # assert_sudo
systemctl restart systemd-journald
test -e /dev/net/tun # check_tun_availability
cloud-init status --wait >/dev/null

# ---------------------------------------------------------------------------- #
# Action
# ---------------------------------------------------------------------------- #

export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y ca-certificates curl gnupg iptables openssl openvpn unbound wget

test -d /etc/openvpn/easy-rsa && rm -rfv /etc/openvpn/easy-rsa/*

# ---------------------------------------------------------------------------- #
# Mapping
# ---------------------------------------------------------------------------- #

find map -type f | while read FILE; do
  sudo mkdir -pv "$(dirname "${FILE/map/}")"
  sudo mv -v "${FILE}" "${FILE/map/}"
done
rm -rfv /root/map

# ---------------------------------------------------------------------------- #
# Ovpn-auth
# ---------------------------------------------------------------------------- #

curl -sSL -o /etc/openvpn/ovpn-auth \
  https://github.com/ufukty/ovpn-auth/releases/download/v1.0.9/ovpn-auth-v1.0.9-linux-amd64
md5sum /etc/openvpn/ovpn-auth | grep "050f02af158934591c1a240e25ca93ab"
chmod 755 /etc/openvpn/ovpn-auth
chown root:root /etc/openvpn/ovpn-auth

# remove_password_change_requirement
# sed --in-place -E 's/root:(.*):0:0:(.*):/root:\1:18770:0:\2:/g' /etc/shadow;
