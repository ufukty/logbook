#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

REGION_SLUG="${1:?}"

sudo -v
sudo openvpn "${STAGE:?}/artifacts/vpn/dth-do-${REGION_SLUG:?}-provisioner.ovpn"
# sleep 1 && sudo killall mDNSResponder{,Helper}
sudo -k
