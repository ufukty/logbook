#!/bin/bash

PS4='\033[33m''\D{%H:%M:%S} provision:${LINENO}:''\033[0m '
set -xeuo pipefail

# ---------------------------------------------------------------------------- #
# Assertions
# ---------------------------------------------------------------------------- #

: "${SSH_PUB_KEYS:?}"
: "${VPS_SUDO_USER_PASSWD_HASH:?}"
: "${VPS_SUDO_USER:?}"

# ---------------------------------------------------------------------------- #
# Wait
# ---------------------------------------------------------------------------- #

timeout 180 bash -c "until stat /var/lib/cloud/instance/boot-finished 2>/dev/null; do sleep 2; done"

# ---------------------------------------------------------------------------- #
# First apt
# ---------------------------------------------------------------------------- #

export DEBIAN_FRONTEND=noninteractive
apt update -y
apt upgrade -y
apt-get install -y tree jq fail2ban openssh-server unattended-upgrades

# ---------------------------------------------------------------------------- #
# Mapping
# ---------------------------------------------------------------------------- #

find map -type f | while read FILE; do
  sudo mkdir -pv "$(dirname "${FILE/map/}")"
  sudo mv -v "${FILE}" "${FILE/map/}"
done
rm -rfv /root/map

# ---------------------------------------------------------------------------- #
# Sudo user
# ---------------------------------------------------------------------------- #

useradd -m -s /bin/bash "${VPS_SUDO_USER}"

mkdir -p "/home/${VPS_SUDO_USER}/.ssh"
chmod 700 "/home/${VPS_SUDO_USER}/.ssh"

echo "${SSH_PUB_KEYS}" >"/home/${VPS_SUDO_USER}/.ssh/authorized_keys"
chmod 600 "/home/${VPS_SUDO_USER}/.ssh/authorized_keys"

ssh-keygen -t "ed25519" -f "/home/${VPS_SUDO_USER}/.ssh/id_ed25519" -C base -N ""
chmod 600 "/home/${VPS_SUDO_USER}/.ssh/id_ed25519"
chmod 644 "/home/${VPS_SUDO_USER}/.ssh/id_ed25519.pub"

chown -R "${VPS_SUDO_USER}:${VPS_SUDO_USER}" "/home/${VPS_SUDO_USER}/.ssh"

echo "${VPS_SUDO_USER} ALL=(ALL) NOPASSWD: ALL" >>/etc/sudoers
echo "AllowUsers ${VPS_SUDO_USER}" >>/etc/ssh/sshd_config

# ---------------------------------------------------------------------------- #
# Root account
# ---------------------------------------------------------------------------- #

rm -v /root/.ssh/authorized_keys
echo "root:${VPS_SUDO_USER_PASSWD_HASH}" | sudo chpasswd --encrypted
sed 's/91/96/' "/root/.bash_profile" >>"/etc/bash.bashrc" # PS4 coloring

# ---------------------------------------------------------------------------- #
# Validate config
# ---------------------------------------------------------------------------- #

systemctl enable fail2ban
systemctl start fail2ban
systemctl status fail2ban
systemctl restart ssh
