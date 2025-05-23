#!/usr/bin/env bash

test "$1" != "-B" && is_up_to_date .completion.timestamp && echo "up to date" && exit 0

PS4="\033[34m\D{%H:%M:%S} ${PWD/"$WORKSPACE/"/}/$0:${LINENO}:\033[0m "
set -xe

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${DO_SSH_KEY_ID}"
: "${STAGE:?}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Git ignored secret files
# ---------------------------------------------------------------------------- #

test -f map/etc/openvpn/ovpn_auth_database.yml

# ---------------------------------------------------------------------------- #
# Values
# ---------------------------------------------------------------------------- #

VPS_HOME="/home/${VPS_SUDO_USER:?}"

BASE_IMAGE_PREFIX="logbook_builder_base"
BASE_IMAGE_ID="$(
  doctl compute image list-user --output json |
    jq -r --arg prefix "$BASE_IMAGE_PREFIX" '[.[] | select(.name | startswith($prefix)) ] | max_by(.created_at).id'
)"
REGION="nyc3"
SIZE="s-1vcpu-1gb"
TRANSFER_REGIONS=(nyc1 nyc2 sfo2 sfo3)

SCM="$(git describe --always --dirty)"
DROPLET_NAME="logbook-builder-vpn-$(date +%y-%m-%d-%H-%M-%S)"
SNAPSHOT_NAME="${DROPLET_NAME//-/_}"
LOG_FILE="logs/$(date +%y.%m.%d.%H.%M.%S)-${SCM}.log"

# ---------------------------------------------------------------------------- #
# Creation
# ---------------------------------------------------------------------------- #

TMP="$(mktemp)"

doctl compute droplet create "${DROPLET_NAME:?}" \
  --image "${BASE_IMAGE_ID:?}" \
  --region "${REGION:?}" \
  --size "${SIZE:?}" \
  --ssh-keys "${DO_SSH_KEY_ID:?}" \
  --enable-private-networking \
  --wait \
  --verbose \
  --output json >"$TMP"

ID="$(jq -r '.[0].id' "$TMP")"
IP="$(jq -r '.[0].networks.v4.[] | select(.type == "public").ip_address' "$TMP")"

: "${IP:?}"

# ---------------------------------------------------------------------------- #
# Set cleanup
# ---------------------------------------------------------------------------- #

cleanup() {
  EC=$?
  if test $EC -eq 0 && test "$ID"; then
    doctl compute droplet delete "$ID" --force
    rm -rv "$TMP"
  else
    echo "Connect to troubleshoot: ssh root@$IP or ssh $VPS_SUDO_USER@$IP"
  fi
  tput bel
  exit $EC
}

trap cleanup EXIT

# ---------------------------------------------------------------------------- #
# Provisioning
# ---------------------------------------------------------------------------- #

ping -o "$IP" && until ssh -i "$STAGE/secrets/ssh/do" "$VPS_SUDO_USER@$IP" exit; do sleep 5; done # wait

scp -i "$STAGE/secrets/ssh/do" provision.sh "$VPS_SUDO_USER@$IP:provision.sh"
rsync -e "ssh -i '$STAGE/secrets/ssh/do'" --verbose --recursive "./map" "$VPS_SUDO_USER@$IP:"

# shellcheck disable=SC2087
ssh -i "$STAGE/secrets/ssh/do" "$VPS_SUDO_USER@$IP" >"$LOG_FILE" 2>&1 <<EOF
  set -e
  cd "$VPS_HOME"
  sudo --preserve-env bash provision.sh
  rm -rf "$VPS_HOME/provision.sh" "$VPS_HOME/map"
  sudo shutdown -h now
EOF

# ---------------------------------------------------------------------------- #
# Snapshot
# ---------------------------------------------------------------------------- #

doctl compute droplet-action snapshot "${ID:?}" --snapshot-name "$SNAPSHOT_NAME" --wait --verbose

# do not use the action id from previous output
SNAPSHOT_ID="$(
  doctl compute snapshot list --output json | jq -r --arg id "$ID" '.[] | select(.resource_id == $id).id'
)"

for TRANSFER_REGION in "${TRANSFER_REGIONS[@]}"; do
  doctl compute image-action transfer "$SNAPSHOT_ID" --region "$TRANSFER_REGION"
done

sleep 100

until test "$(doctl compute snapshot list --output json | jq --arg id "$ID" '.[] | select(.resource_id == $id).regions | length')" -eq "$((${#TRANSFER_REGIONS[@]} + 1))"; do
  sleep 10
done

doctl compute snapshot list --output json | jq --arg id "$ID" '.[] | select(.resource_id == $id)'
