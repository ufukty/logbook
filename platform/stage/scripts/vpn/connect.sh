#!/usr/local/bin/bash

REGION_SLUG="$1" && shift

sudo -v
sudo openvpn "${STAGE:?}/artifacts/vpn/dth-do-${REGION_SLUG:?}-provisioner.ovpn"
# sleep 1 && sudo killall mDNSResponder{,Helper}
sudo -k
