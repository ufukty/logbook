#!/usr/local/bin/bash

test "$1" != "-B" && is_up_to_date .completion.timestamp && echo "up to date" && exit 0

PS4='\033[31m''\D{%H:%M:%S} build:${LINENO}:''\033[0m '
set -xeuo pipefail

# ---------------------------------------------------------------------------- #
# Values
# ---------------------------------------------------------------------------- #

VPS_SUDO_USER="olwgtzjzhnvexhpr"

BASE="ubuntu-22-04-x64"
REGION="fra1"
SIZE="s-1vcpu-1gb"
SSH_KEY_IDs="41814107"

TRANSFER_REGIONS=() # ("nyc3" "ams3")

FOLDER="$(basename "$PWD")"
DROPLET_NAME="builder-${FOLDER:?}-$(date +%y-%m-%d-%H-%M-%S)"
SNAPSHOT_NAME="build_${FOLDER:?}_$(date +%y_%m_%d_%H_%M_%S)"

# ---------------------------------------------------------------------------- #
# Creation
# ---------------------------------------------------------------------------- #

DROPLET="$(doctl compute droplet create "${DROPLET_NAME:?}" --image "${BASE:?}" --region "${REGION:?}" --size "${SIZE:?}" --ssh-keys "${SSH_KEY_IDs:?}" --tag-name "${FOLDER:?}" --wait --verbose --no-header)"
ID="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $1 }')"
IP="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $3 }')"

cleanup() {
  EC=$?
  test "$ID" && test "$1" != "-d" && doctl compute droplet delete "$ID" --force
  test $EC -eq 0 && touch .completion.timestamp
  tput bel
  exit $EC
}

trap cleanup EXIT

# ---------------------------------------------------------------------------- #
# Provisioning
# ---------------------------------------------------------------------------- #

ping -o "${IP:?}" && until ssh "root@${IP:?}" exit; do sleep 5; done # wait

rsync --verbose --recursive -e ssh "./map" "root@${IP:?}:/root/"

export ANSIBLE_CONFIG="ansible/ansible.cfg"
ansible-playbook -i "${IP:?}," -u root ansible/playbook.yml
ssh "${VPS_SUDO_USER:?}@${IP:?}" sudo shutdown -h now

# ---------------------------------------------------------------------------- #
# Snapshot
# ---------------------------------------------------------------------------- #

doctl compute droplet-action snapshot "${ID:?}" --snapshot-name "${SNAPSHOT_NAME:?}" --wait --verbose

SNAPSHOT_ID="$(doctl compute snapshot list | grep "$ID" | awk '{ print $1 }')" # do not use the action id from previous output

for TRANSFER_REGION in "${TRANSFER_REGIONS[@]}"; do
  doctl compute image-action transfer "$SNAPSHOT_ID" --region "$TRANSFER_REGION" --wait
done

doctl compute snapshot list | grep -e "$SNAPSHOT_ID" -e "Created at"
