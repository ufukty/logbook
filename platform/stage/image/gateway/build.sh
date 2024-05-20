#!/usr/local/bin/bash

test "$1" != "-B" && is_up_to_date .completion.timestamp && echo "up to date" && exit 0

PS4="\033[32m/\$(basename \"\${BASH_SOURCE}\"):\${LINENO}\033[0m\033[33m\${FUNCNAME[0]:+/\${FUNCNAME[0]}():}\033[0m "
set -x
set -v
set -e

VPS_SUDO_USER="olwgtzjzhnvexhpr"
VPS_HOME="/home/${VPS_SUDO_USER:?}"
IPTABLES_PUBLIC_ETHERNET_INTERFACE="eth0"
IPTABLES_PRIVATE_ETHERNET_INTERFACE="eth1"

BASE_IMAGE_PREFIX="build_internal"
BASE_IMAGE_ID="$(doctl compute image list-user --format Name,ID --no-header | grep "^$BASE_IMAGE_PREFIX" | sort | tail -n 1 | awk '{ print $2 }')"
REGION="fra1"
SIZE="s-1vcpu-1gb"
SSH_KEY_IDs="41814107"
VPC_UUID="$(doctl vpcs list | grep logbook-fra1 | awk '{ print $1 }')"

TRANSFER_REGIONS=("nyc3" "ams3")

FOLDER="$(basename "$PWD")"
DROPLET_NAME="builder-${FOLDER:?}-$(date +%y-%m-%d-%H-%M-%S)"
SNAPSHOT_NAME="build_${FOLDER:?}_$(date +%y_%m_%d_%H_%M_%S)"

DROPLET="$(doctl compute droplet create "${DROPLET_NAME:?}" --image "${BASE_IMAGE_ID:?}" --region "${REGION:?}" --size "${SIZE:?}" --ssh-keys "${SSH_KEY_IDs:?}" --tag-name "${FOLDER:?}" --vpc-uuid "${VPC_UUID:?}" --enable-private-networking --wait --verbose --no-header --format ID,PrivateIPv4)"
ID="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $1 }')"
IP="$(echo "$DROPLET" | tail -n 1 | awk '{ print  $2 }')"

cleanup() {
    EC=$?
    test $EC -eq 0 && test "$ID" && doctl compute droplet delete "$ID" --force
    tput bel
    exit $EC
}

trap cleanup EXIT

ping -o "${IP:?}" && until ssh "${VPS_SUDO_USER:?}@${IP:?}" exit; do sleep 5; done # wait

rsync --verbose --recursive -e ssh "./files" "${VPS_SUDO_USER:?}@${IP:?}:${VPS_HOME:?}/"
ssh "${VPS_SUDO_USER:?}@$IP" bash <<EOF
    set -e
    set -v
    cd "${VPS_HOME:?}/files"
    sudo --preserve-env \
        IPTABLES_PUBLIC_ETHERNET_INTERFACE="${IPTABLES_PUBLIC_ETHERNET_INTERFACE:?}" \
        IPTABLES_PRIVATE_ETHERNET_INTERFACE="${IPTABLES_PRIVATE_ETHERNET_INTERFACE:?}" \
        bash golden-image.sh
    cd "${VPS_HOME:?}"
    rm -rf "${VPS_HOME:?}/files"
    sudo shutdown -h now
EOF

doctl compute droplet-action snapshot "${ID:?}" --snapshot-name "${SNAPSHOT_NAME:?}" --wait --verbose
SNAPSHOT_ID="$(doctl compute snapshot list | grep "$ID" | awk '{ print $1 }')" # do not use the action id from previous output

for TRANSFER_REGION in "${TRANSFER_REGIONS[@]}"; do
    doctl compute image-action transfer "$SNAPSHOT_ID" --region "$TRANSFER_REGION" --wait --verbose
done

doctl compute snapshot list | grep -e "$SNAPSHOT_ID" -e "Created at"

touch .completion.timestamp
