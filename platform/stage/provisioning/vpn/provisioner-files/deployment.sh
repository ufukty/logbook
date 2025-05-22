#!/bin/bash

# ---------------------------------------------------------------------------- #
# Documentation
# ---------------------------------------------------------------------------- #

# Summary:           This script creates a OpenVPN server.

# Features:        - It Installs EasyRSA & OpenVPN, configures
#                    OpenVPN & iptables, creates .ovpn files for
#                    clients.
#                  - Clients can access to internet through
#                    server.
#                  - Clients can access to hosts available in
#                    the LAN of server
#                  - Clients are isolated from each other.
#                  - Same settings can be changed with optional
#                    variables including the encryption.

# Requirements:    - Fresh Ubuntu 20.04 server
#                  - Must run with sudo or by root user

# Example usage:     PUBLIC_IP="192.168.33.8" \
#                    VPC_ADDRESS="10.170.0.0" \
#                    PUBLIC_ETHERNET_INTERFACE="eth0" \
#                    PRIVATE_ETHERNET_INTERFACE="eth1" \
#                    SERVER_NAME="my_server" \
#                    SERVER_NAME="my_server_common_name" \
#                    sudo --preserve-env bash sh.sh \
#                    my_client_1 my_client_2 my_client_n

# Fork:              This script is derived from another script
#                    https://github.com/angristan/openvpn-install

PS4="\033[32m/\$(basename \"\${BASH_SOURCE}\"):\${LINENO}\033[0m\033[33m\${FUNCNAME[0]:+/\${FUNCNAME[0]}():}\033[0m "

set -x
set -v
set -e
set -E

# ---------------------------------------------------------------------------- #
# Required Environment Variables
# ---------------------------------------------------------------------------- #

: "${USER_ACCOUNT_NAME:?}"

: "${PUBLIC_IP:?}"
: "${PRIVATE_IP:?}"

# OpenVPN to handle underlying networking jobs. e.g. $OPENVPN_SUBNET_ADDRESS
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

# Those will be used as login credentials to ovpn-auth
: "${OVPN_AUTH_USERNAME:?}"
: "${OVPN_AUTH_HASH:?}"
: "${OVPN_AUTH_TOTP:?}"

# ---------------------------------------------------------------------------- #
# Optional Environment Variables
# ---------------------------------------------------------------------------- #

# OPENVPN_PROTOCOL valid values:
# [ udp, tcp ]
OPENVPN_PROTOCOL="${OPENVPN_PROTOCOL:-"tcp"}"
OPENVPN_PORT="${OPENVPN_PORT:-"443"}"

# ENCRYPTION_TLS_SIG valid values:
# [ tls-crypt, tls-auth ]
ENCRYPTION_TLS_SIG=${ENCRYPTION_TLS_SIG:-"tls-crypt"}

# ENCRYPTION_CIPHER valid values:
# [ AES-128-GCM, AES-192-GCM, AES-256-GCM, AES-128-CBC, AES-192-CBC, AES-256-CBC ]
ENCRYPTION_CIPHER="${ENCRYPTION_CIPHER:-"AES-128-GCM"}"

# ENCRYPTION_CERT_TYPE valid values:
# [ ECDSA, RSA ]
ENCRYPTION_CERT_TYPE="${ENCRYPTION_CERT_TYPE:-"ECDSA"}"

# !!! Only usable when ENCRYPTION_CERT_TYPE="ECDSA"
# ENCRYPTION_CERT_CURVE valid values:
# [ prime256v1, secp384r1, secp521r1 ]
ENCRYPTION_ECDSA_CERT_CURVE="${ENCRYPTION_ECDSA_CERT_CURVE:-"prime256v1"}"

# !!! Only usable when ENCRYPTION_CERT_TYPE="ECDSA"
# ENCRYPTION_ECDSA_CC_CIPHER valid values:
# [ ECDHE-ECDSA-AES-128-GCM-SHA256, ECDHE-ECDSA-AES-256-GCM-SHA384 ]
ENCRYPTION_ECDSA_CC_CIPHER="${ENCRYPTION_ECDSA_CC_CIPHER:-"TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256"}"

# !!! Only usable when ENCRYPTION_CERT_TYPE="RSA"
# RSA_KEY_SIZE valid values:
# [ 2048, 3072, 4096 ]
ENCRYPTION_RSA_KEY_SIZE="${ENCRYPTION_RSA_KEY_SIZE:-"2048"}"

# !!! Only usable when ENCRYPTION_CERT_TYPE="RSA"
# ENCRYPTION_RSA_CC_CIPHER valid values:
# [ ECDHE-RSA-AES-128-GCM-SHA256, ECDHE-RSA-AES-256-GCM-SHA384 ]
ENCRYPTION_RSA_CC_CIPHER="${ENCRYPTION_RSA_CC_CIPHER:-"ECDHE-RSA-AES-128-GCM-SHA256"}"

# ENCRYPTION_DH_TYPE valid values:
# [ ECDH, DH ]
ENCRYPTION_DH_TYPE="${ENCRYPTION_DH_TYPE:-"ECDH"}"

# !!! Only usable when ENCRYPTION_DH_TYPE="ECDH"
# ENCRYPTION_ECDH_DH_CURVE valid values:
# [ prime256v1, secp384r1, secp521r1 ]
ENCRYPTION_ECDH_CURVE="${ENCRYPTION_ECDH_CURVE:-"prime256v1"}"

# !!! Only usable when ENCRYPTION_DH_TYPE="DH"
# ENCRYPTION_DH_KEY_SIZE valid values:
# [ 2048, 3072, 4096 ]
ENCRYPTION_DH_KEY_SIZE="${ENCRYPTION_DH_KEY_SIZE:-"2048"}"

# ENCRYPTION_HMAC_ALG valid values:
# [ SHA-256, SHA-384, SHA-512 ]
# - When GCM type ciphers are used, the algorithm is used only for
#   encryption of tls-auth packets from the control channel.
# - If, CBC type ciphers are used, the algorithm is used in addition
#   for authenticates data channel packets too.
ENCRYPTION_HMAC_ALG="${ENCRYPTION_HMAC_ALG:-"SHA256"}"

# TLS_SIG valid values:
# [ tls-crypt, tls-auth ]
# - Those will add additional layer of security to the control channel.
# - tls-auth authenticates the packets, while tls-crypt authenticate
#   and encrypt them.
TLS_SIG="${TLS_SIG:-"tls-crypt"}"

# ---------------------------------------------------------------------------- #
# Commons
# ---------------------------------------------------------------------------- #

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
  chmod 700 -R "${PROVISIONER_FILES:?}/map"
  chown root:root -R "${PROVISIONER_FILES:?}/map"
  rsync --verbose --recursive "${PROVISIONER_FILES:?}/map/" "/"
  rm -rfv "${PROVISIONER_FILES:?}/map"
}

export DEBIAN_FRONTEND=noninteractive

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

EASYRSA_CA_NAME="$SERVER_NAME-certificate-authority"
EASYRSA_SERVER_NAME="$SERVER_NAME-server"
# OVPN_GENERATED_TOTP_SHARE_FILE="/home/${USER_ACCOUNT_NAME:?}/artifacts/totp-share.txt"

mkdir -p "/etc/openvpn/easy-rsa/generated"
mkdir -p "/home/${USER_ACCOUNT_NAME:?}/artifacts"
echo -n "${EASYRSA_CA_NAME:?}" >"/etc/openvpn/easy-rsa/generated/ca_name"
echo -n "${EASYRSA_SERVER_NAME:?}" >"/etc/openvpn/easy-rsa/generated/server_name"

function configure_easy_rsa() (
  cd /etc/openvpn/easy-rsa

  if [[ $ENCRYPTION_CERT_TYPE == "ECDSA" ]]; then
    echo "set_var EASYRSA_ALGO           \"ec\"" >>vars
    echo "set_var EASYRSA_CURVE          \"$ENCRYPTION_ECDSA_CERT_CURVE\"" >>vars
  elif [[ $ENCRYPTION_CERT_TYPE == "RSA" ]]; then
    echo "set_var EASYRSA_KEY_SIZE       \"$ENCRYPTION_RSA_KEY_SIZE\"" >>vars
  fi

  # Create the PKI, set up the CA, the DH params and the server certificate
  ./easyrsa --vars=./vars init-pki
  ./easyrsa --vars=./vars --batch --req-cn="${EASYRSA_CA_NAME:?}" build-ca nopass

  # ECDH keys are generated on-the-fly so we don't need to generate them beforehand
  if [[ $ENCRYPTION_DH_TYPE == "DH" ]]; then
    openssl dhparam -out dh.pem "$ENCRYPTION_DH_KEY_SIZE"
  fi

  ./easyrsa --vars=./vars --batch build-server-full "${EASYRSA_SERVER_NAME:?}" nopass
  EASYRSA_CRL_DAYS=3650 ./easyrsa --vars=./vars gen-crl

  # Generate tls-crypt or tls-auth key
  openvpn --genkey secret "/etc/openvpn/${TLS_SIG:?}.key"
)

function configure_openvpn() (
  cd /etc/openvpn/easy-rsa

  # Move all the generated files
  cp \
    "/etc/openvpn/easy-rsa/pki/ca.crt" \
    "/etc/openvpn/easy-rsa/pki/private/ca.key" \
    "/etc/openvpn/easy-rsa/pki/issued/${EASYRSA_SERVER_NAME:?}.crt" \
    "/etc/openvpn/easy-rsa/pki/private/${EASYRSA_SERVER_NAME:?}.key" \
    "/etc/openvpn/easy-rsa/pki/crl.pem" \
    "/etc/openvpn"

  if [[ $ENCRYPTION_DH_TYPE == "ECDH" ]]; then
    DH_CONF_STR="dh none \necdh-curve $ENCRYPTION_ECDH_CURVE"
  elif [[ $ENCRYPTION_DH_TYPE == "DH" ]]; then
    cp dh.pem /etc/openvpn
    DH_CONF_STR="dh dh.pem"
  fi

  chmod 644 /etc/openvpn/crl.pem

  # Find out if the machine uses nogroup or nobody for the permissionless group
  if grep -qs "^nogroup:" /etc/group; then
    NOGROUP=nogroup
  else
    NOGROUP=nobody
  fi

  if [[ $ENCRYPTION_CERT_TYPE == "ECDSA" ]]; then
    ENCRYPTION_CC_CIPHER="$ENCRYPTION_ECDSA_CC_CIPHER"
  elif [[ $ENCRYPTION_CERT_TYPE == "RSA" ]]; then
    ENCRYPTION_CC_CIPHER="$ENCRYPTION_RSA_CC_CIPHER"
  fi

  # "Populating the configure file at: /etc/openvpn/server.conf"
  sed --in-place \
    -e "s;{{DH_CONF_STR}};${DH_CONF_STR:?};g" \
    -e "s;{{ENCRYPTION_CC_CIPHER}};${ENCRYPTION_CC_CIPHER:?};g" \
    -e "s;{{ENCRYPTION_CIPHER}};${ENCRYPTION_CIPHER:?};g" \
    -e "s;{{ENCRYPTION_HMAC_ALG}};${ENCRYPTION_HMAC_ALG:?};g" \
    -e "s;{{NOGROUP}};${NOGROUP:?};g" \
    -e "s;{{OPENVPN_PORT}};${OPENVPN_PORT:?};g" \
    -e "s;{{OPENVPN_PROTOCOL}};${OPENVPN_PROTOCOL:?};g" \
    -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS:?};g" \
    -e "s;{{OPENVPN_SUBNET_MASK}};${OPENVPN_SUBNET_MASK:?};g" \
    -e "s;{{EASYRSA_SERVER_NAME}};${EASYRSA_SERVER_NAME:?};g" \
    -e "s;{{TLS_SIG}};${TLS_SIG:?};g" \
    -e "s;{{UNBOUND_ADDRESS}};${UNBOUND_ADDRESS:?};g" \
    -e "s;{{VPC_RANGE_ADDRESS}};${VPC_RANGE_ADDRESS:?};g" \
    -e "s;{{VPC_RANGE_MASK}};${VPC_RANGE_MASK:?};g" \
    /etc/openvpn/server.conf

  mkdir -p /etc/openvpn/ccd # Create client-config-dir dir
  mkdir -p /var/log/openvpn # Create log dir

  sysctl --system # "Apply sysctl rules"

  # If SELinux is enabled and a custom port was selected, we need this
  if hash sestatus 2>/dev/null; then
    if sestatus | grep "Current mode" | grep -qs "enforcing"; then
      if [[ ${OPENVPN_PORT:?} != '1194' ]]; then
        semanage port -a -t openvpn_port_t -p "${OPENVPN_PROTOCOL:?}" "${OPENVPN_PORT:?}"
      fi
    fi
  fi

  systemctl enable openvpn
  systemctl start openvpn
)

function configure_ovpn_auth() {
  sed --in-place \
    -e "s;<<OVPN_AUTH_USERNAME>>;${OVPN_AUTH_USERNAME:?};g" \
    -e "s;<<OVPN_AUTH_HASH>>;${OVPN_AUTH_HASH:?};g" \
    -e "s;<<OVPN_AUTH_TOTP>>;${OVPN_AUTH_TOTP:?};g" \
    /etc/openvpn/ovpn_auth_database.yml

  chmod 744 /etc/openvpn/ovpn_auth_database.yml
  chown root:root /etc/openvpn/ovpn_auth_database.yml
}

function configure_iptables() {
  sed --in-place \
    -e "s;{{PRIVATE_ETHERNET_INTERFACE}};${PRIVATE_ETHERNET_INTERFACE:?};g" \
    -e "s;{{PUBLIC_ETHERNET_INTERFACE}};${PUBLIC_ETHERNET_INTERFACE:?};g" \
    -e "s;{{OPENVPN_PROTOCOL}};${OPENVPN_PROTOCOL:?};g" \
    -e "s;{{OPENVPN_PORT}};${OPENVPN_PORT:?};g" \
    -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS:?};g" \
    /etc/iptables/iptables-rules.v4

  chmod 644 /etc/iptables/iptables-rules.v4
  chmod 644 /etc/systemd/system/iptables-activation.service

  systemctl daemon-reload
  systemctl enable iptables-activation
}

function configure_unbound() {
  sed --in-place \
    -e "s;{{UNBOUND_ADDRESS}};${UNBOUND_ADDRESS:?};g" \
    -e "s;{{OPENVPN_SUBNET_ADDRESS}};${OPENVPN_SUBNET_ADDRESS:?};g" \
    /etc/unbound/unbound.conf.d/unbound.conf
  # -e "s;{{HOST_ADDRESS}};$PRIVATE_IP;g" \
  # -e "s;{{VPC_CIDR}};$VPC_CIDR;g" \

  systemctl enable unbound
  systemctl restart unbound
}

assert_sudo
restart_journald
check_tun_availability
wait_cloud_init

configure_easy_rsa
configure_openvpn
configure_ovpn_auth
configure_iptables
configure_unbound
