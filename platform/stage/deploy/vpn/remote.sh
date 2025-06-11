#!/usr/bin/env bash

PS4='\033[32m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Wait
# ---------------------------------------------------------------------------- #

export DEBIAN_FRONTEND=noninteractive
test $EUID -eq 0
systemctl restart systemd-journald
test -e /dev/net/tun
cloud-init status --wait >/dev/null

# ---------------------------------------------------------------------------- #
# Configure OpenVPN
# ---------------------------------------------------------------------------- #

cd "$HOME"
mv crl.pem ovpn_auth_database.yml tls-crypt.key /etc/openvpn/

cd /etc/openvpn
chown -R openvpn:openvpn *
chmod 600 server.key ovpn_auth_database.yml
chmod 640 ./{crl.pem,tls-crypt.key,server.conf}
mkdir -p /var/log/openvpn
sysctl --system
systemctl enable openvpn
systemctl restart openvpn

# ---------------------------------------------------------------------------- #
# Configure iptables
# ---------------------------------------------------------------------------- #

chmod 644 /etc/iptables/iptables-rules.v4
chmod 644 /etc/systemd/system/iptables-activation.service
systemctl daemon-reload
systemctl enable iptables-activation

# ---------------------------------------------------------------------------- #
# Configure unbound
# ---------------------------------------------------------------------------- #

systemctl enable unbound
systemctl restart unbound
