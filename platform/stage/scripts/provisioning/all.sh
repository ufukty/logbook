#!/usr/bin/env bash

PS4='\033[31m$0:$LINENO:\033[0m '
set -xe

bash scripts/provision/vpn-up --auto-approve
bash scripts/provision/app-up --auto-approve
bash scripts/digitalocean/ssh-keyscan
bash scripts/deploy
