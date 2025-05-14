#!/usr/local/bin/bash

test "$1" != "-B" && is_up_to_date .completion.timestamp && echo "up to date" && exit 0

PS4='\033[31m''\D{%H:%M:%S} build:${LINENO}:''\033[0m '
set -xeuo pipefail

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${DO_SSH_FINGERPRINT}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Values
# ---------------------------------------------------------------------------- #

BASE="${BASE:-"ubuntu-24-04-x64"}"
REGION="nyc3"
SIZE="${SIZE:-"s-1vcpu-1gb"}"

TRANSFER_REGIONS=() # ("nyc3" "ams3")

FOLDER="$(basename "$PWD")"
DROPLET_NAME="builder-${FOLDER:?}-$(date +%y-%m-%d-%H-%M-%S)"
SNAPSHOT_NAME="build_${FOLDER:?}_$(date +%y_%m_%d_%H_%M_%S)"

# ---------------------------------------------------------------------------- #
# Creation
# ---------------------------------------------------------------------------- #

mkdir -p tmp

doctl compute droplet create "${DROPLET_NAME:?}" \
  --image "${BASE:?}" \
  --region "${REGION:?}" \
  --size "${SIZE:?}" \
  --ssh-keys "${SSH_KEY_FINGERPRINT:?}" \
  --wait \
  --verbose \
  --output json >tmp/droplet.json

ID="$(jq -r '.[0].id' tmp/droplet.json)"
IP="$(jq -r '.[0].networks.v4.[] | select(.type == "public").ip_address' tmp/droplet.json)"

: "${IP:?}"

# ---------------------------------------------------------------------------- #
# Set cleanup
# ---------------------------------------------------------------------------- #

cleanup() {
  EC=$?
  if test $EC -eq 0 && test "$ID"; then
    doctl compute droplet delete "$ID" --force
    rm -rv tmp
  else
    echo "Connect to troubleshoot: ssh root@${IP} or ssh ${VPS_SUDO_USER}@${IP}"
  fi
  tput bel
  exit $EC
}

trap cleanup EXIT

# ---------------------------------------------------------------------------- #
# Provisioning
# ---------------------------------------------------------------------------- #

ping -o "${IP}" && until ssh "root@${IP}" exit; do sleep 5; done # wait

scp provision.sh "root@${IP}:provision.sh"
rsync --verbose --recursive -e ssh "./map" "root@${IP}:"

mkdir -p logs

ssh "root@${IP}" >"$LOG_FILE" 2>&1 <<EOF
  set -e
  SSH_PUB_KEYS="${SSH_PUB_KEYS}" \
  VPS_SUDO_USER_PASSWD_HASH="${VPS_SUDO_USER_PASSWD_HASH}" \
  VPS_SUDO_USER="${VPS_SUDO_USER}" \
  bash provision.sh
  rm provision.sh
  sudo shutdown -h now
EOF

# ---------------------------------------------------------------------------- #
# Snapshot
# ---------------------------------------------------------------------------- #

doctl compute droplet-action snapshot "${ID:?}" --snapshot-name "${SNAPSHOT_NAME:?}" --wait --verbose

# SNAPSHOT_ID="$(doctl compute snapshot list | grep "$ID" | awk '{ print $1 }')" # do not use the action id from previous output

# for TRANSFER_REGION in "${TRANSFER_REGIONS[@]}"; do
#   doctl compute image-action transfer "$SNAPSHOT_ID" --region "$TRANSFER_REGION" --wait
# done

# doctl compute snapshot list | grep -e "$SNAPSHOT_ID" -e "Created at"
