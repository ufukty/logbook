#!/usr/local/bin/bash

PS4="\033[36m$(grealpath --relative-to="$(dirname "$WORKSPACE")" "$(pwd -P)")\033[32m/\$(basename \"\${BASH_SOURCE}\"):\${LINENO}\033[0m\033[33m\${FUNCNAME[0]:+/\${FUNCNAME[0]}():}\033[0m "
set -x
set -e

BASE="ubuntu-22-04-x64"
REGION="fra1"
SIZE="s-1vcpu-1gb"
SSH_KEY_IDs="41814107"

FOLDER="$(basename "$PWD")"
DROPLET_NAME="builder-${FOLDER:?}-$(date +%y-%m-%d-%H-%M-%S)"
SNAPSHOT_NAME="build_${FOLDER:?}_$(date +%y_%m_%d_%H_%M_%S)"

DROPLET="$(
    doctl compute droplet create "${DROPLET_NAME:?}" \
        --image "${BASE:?}" \
        --region "${REGION:?}" \
        --size "${SIZE:?}" \
        --ssh-keys "${SSH_KEY_IDs:?}" \
        --tag-name "${FOLDER:?}" \
        --wait \
        --verbose \
        --no-header
)"
ID="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $1 }')"
IP="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $3 }')"

cleanup() {
    EC=$?
    test $EC -eq 0 && test "$ID" && doctl compute droplet delete "$ID" --force
    exit $EC
}

trap cleanup EXIT

ping -o "${IP:?}" && until ssh "root@${IP:?}" exit; do sleep 5; done # wait

rsync --verbose --recursive -e ssh "./map" "root@${IP:?}:/root/"

export ANSIBLE_CONFIG="ansible/ansible.cfg"
ansible-playbook -i "${IP:?}," -u root ansible/playbook.yml

# VPS_USERNAME="olwgtzjzhnvexhpr"
# VPS_HOME="/home/${VPS_USERNAME:?}"

doctl compute droplet-action snapshot "${ID:?}" \
    --snapshot-name "${SNAPSHOT_NAME:?}" \
    --wait \
    --verbose

echo "$PWD"

