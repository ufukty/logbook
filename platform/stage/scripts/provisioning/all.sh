#!/usr/local/bin/bash

bash scripts/provision/vpn-up --auto-approve
bash scripts/provision/app-up --auto-approve
bash scripts/digitalocean/ssh-keyscan
bash scripts/deploy
