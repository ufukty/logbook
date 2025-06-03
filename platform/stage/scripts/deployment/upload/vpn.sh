#!/usr/bin/env bash

# Example usage:
#
# VPC_ADDRESS="10.170.0.0" \
# PUBLIC_ETHERNET_INTERFACE="eth0" \
# PRIVATE_ETHERNET_INTERFACE="eth1" \
# SERVER_NAME="my_server" \
# SERVER_NAME="my_server_common_name" \
# sudo --preserve-env bash sh.sh \
# my_client_1 my_client_2 my_client_n
#
# Prior art:
#
# The steps in this script follows:
# https://github.com/angristan/openvpn-install

PS4='\033[32m$0:$LINENO\033[0m: '
set -xe

# ---------------------------------------------------------------------------- #
# Required Environment Variables
# ---------------------------------------------------------------------------- #

: "${VPS_SUDO_USER:?}"

: "${PRIVATE_IP:?}"

# OpenVPN to handle underlying networking jobs.
: "${OPENVPN_SUBNET_ADDRESS:?}"
: "${OPENVPN_SUBNET_MASK:?}"

# PUBLIC_ETHERNET_INTERFACE
# Like eth0
: "${PUBLIC_ETHERNET_INTERFACE:?}"

# The name of the ethernet adapter.
# Used for connecting public internet. Like eth1
: "${PRIVATE_ETHERNET_INTERFACE:?}"

# Used for EasyRSA and ovpn TOTP URI.
# It could be an arbitrary string that is unique to each region/provider.
: "${SERVER_NAME:?}"

# ---------------------------------------------------------------------------- #
# Runtime Variables
# ---------------------------------------------------------------------------- #

UNBOUND_ADDRESS="$(ip -json route list dev "${PRIVATE_ETHERNET_INTERFACE:?}" | jq -r '.[0].prefsrc')" # IP points to itself
VPC_CIDR="$(ip -json route list dev "${PRIVATE_ETHERNET_INTERFACE:?}" | jq -r '.[0].dst')"
VPC_RANGE_ADDRESS="$(ipcalc "${VPC_CIDR:?}" --nobinary --nocolor | grep Address | awk '{ print $2 }')"
VPC_RANGE_MASK="$(ipcalc "${VPC_CIDR:?}" --nobinary --nocolor | grep Netmask | awk '{ print $2 }')"

# ---------------------------------------------------------------------------- #
# Definitions
# ---------------------------------------------------------------------------- #

EASYRSA_SERVER_NAME="$SERVER_NAME-server"

# ---------------------------------------------------------------------------- #
# Prep
# ---------------------------------------------------------------------------- #

export DEBIAN_FRONTEND=noninteractive

# Assert Sudo
test $EUID -eq 0

# Restart Journald
systemctl restart systemd-journald

# Check Tun Availability
test -e /dev/net/tun

# Wait Cloud Init
cloud-init status --wait >/dev/null

# ---------------------------------------------------------------------------- #
# Configure OpenVPN
# ---------------------------------------------------------------------------- #

cd /etc/openvpn/easy-rsa

# "Populating the configure file at: /etc/openvpn/server.conf"
sed --in-place \
  -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS};g" \
  -e "s;{{OPENVPN_SUBNET_MASK}};${OPENVPN_SUBNET_MASK};g" \
  -e "s;{{EASYRSA_SERVER_NAME}};${EASYRSA_SERVER_NAME};g" \
  -e "s;{{UNBOUND_ADDRESS}};${UNBOUND_ADDRESS};g" \
  -e "s;{{VPC_RANGE_ADDRESS}};${VPC_RANGE_ADDRESS};g" \
  -e "s;{{VPC_RANGE_MASK}};${VPC_RANGE_MASK};g" \
  /etc/openvpn/server.conf

mkdir -p /etc/openvpn/ccd # Create client-config-dir dir
mkdir -p /var/log/openvpn # Create log dir

sysctl --system # "Apply sysctl rules"

systemctl enable openvpn
systemctl start openvpn

# ---------------------------------------------------------------------------- #
# Configure_iptables
# ---------------------------------------------------------------------------- #

sed --in-place \
  -e "s;{{PRIVATE_ETHERNET_INTERFACE}};${PRIVATE_ETHERNET_INTERFACE:?};g" \
  -e "s;{{PUBLIC_ETHERNET_INTERFACE}};${PUBLIC_ETHERNET_INTERFACE:?};g" \
  -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS:?};g" \
  /etc/iptables/iptables-rules.v4

chmod 644 /etc/iptables/iptables-rules.v4
chmod 644 /etc/systemd/system/iptables-activation.service

systemctl daemon-reload
systemctl enable iptables-activation

# ---------------------------------------------------------------------------- #
# Configure_unbound
# ---------------------------------------------------------------------------- #

sed --in-place \
  -e "s;{{UNBOUND_ADDRESS}};${UNBOUND_ADDRESS:?};g" \
  -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS:?};g" \
  /etc/unbound/unbound.conf.d/unbound.conf
# -e "s;{{HOST_ADDRESS}};$PRIVATE_IP;g" \
# -e "s;{{VPC_CIDR}};$VPC_CIDR;g" \

systemctl enable unbound
systemctl restart unbound
