#!/usr/local/bin/bash

STAGE="${WORKSPACE:?}/platform/stage"
export STAGE

# MARK: SSH overload

alias ssh="ssh -F ${STAGE:?}/artifacts/ssh.conf"
alias scp="scp -F ${STAGE:?}/artifacts/ssh.conf"
_check_env_vars() {
    : "${DIGITALOCEAN_ACCESS_TOKEN:?}"
    : "${TF_VAR_DIGITALOCEAN_TOKEN:?}"
    : "${TF_VAR_OVPN_AUTH_USERNAME:?}"
    : "${TF_VAR_OVPN_AUTH_HASH:?}"
    : "${TF_VAR_OVPN_AUTH_TOTP:?}"
    : "${ANSIBLE_SUDO_USER_PASSWD_HASH:?}"
}
_check_env_vars

_ssh_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD - 1]}"
    opts=$(grep '^Host' ${STAGE:?}/artifacts/ssh.conf 2>/dev/null | grep -v '[?*]' | cut -d ' ' -f 2-)

    COMPREPLY=($(compgen -W "$opts" -- ${cur}))
    return 0
}
complete -F _ssh_completion ssh

# MARK: Utilities

PING_URL="stage.logbook.balaasad.com:8080/api/v1.0.0/ping"

aggregate-ssh-conf() {
    cat "${STAGE:?}/artifacts/ssh.conf.d/"* >"${STAGE:?}/artifacts/ssh.conf"
}

ssh-key-update() {
    touch "${STAGE:?}/artifacts/deployment/service_discovery.json"
    ADDRESSES="$(cat "${STAGE:?}/artifacts/deployment/service_discovery.json" | jq -r '.digitalocean.fra1.services[] | .[] | .ipv4_address_private')"
    echo "$ADDRESSES" | while read ADDRESS; do
        ssh-keygen -R "$ADDRESS" >/dev/null 2>&1
        # ssh-keyscan "$ADDRESS" >>~/.ssh/known_hosts 2>/dev/null
    done
}

update-dns-records() (
    set -x -T -v -e -E
    set -o pipefail
    local GATEWAY_IP
    GATEWAY_IP="$(jq -r '.digitalocean.fra1.services["api-gateway"][0].ipv4_address' <"${STAGE:?}/artifacts/deployment/service_discovery.json")"
    ssh -t fra1-vpn "sudo bash -c 'sed \"s;{{GATEWAY_IP}};${GATEWAY_IP:?};g\" /etc/unbound/unbound.conf.tmpl.d/custom.conf > /etc/unbound/unbound.conf.d/custom.conf && systemctl restart unbound'"
    sudo killall mDNSResponder{,Helper}
)

# MARK: Provision

vpc-up() (
    cd "${STAGE:?}/provisioning/vpc"
    terraform apply "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
)

vpn-up() (
    set -e
    cd "${STAGE:?}/provisioning/vpn"
    terraform apply --auto-approve --var-file="${STAGE:?}/provisioning/vars.tfvars"
    aggregate-ssh-conf
    note "Connect vpn in separate tab [Enter]"
    read # wait
    ssh-key-update
    update-dns-records
)

vpn-down() (
    cd "${STAGE:?}/provisioning/vpn"
    terraform destroy "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
    aggregate-ssh-conf
)

app-up() (
    cd "${STAGE:?}/provisioning/application"
    terraform apply "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
    aggregate-ssh-conf
    ssh-key-update
    update-dns-records
)

app-down() (
    cd "${STAGE:?}/provisioning/application"
    terraform destroy "$@" --var-file="${STAGE:?}/provisioning/vars.tfvars"
    aggregate-ssh-conf
)

# MARK: Deployment

deploy() (
    export PROGRAM_NAME="$1" && shift
    cd "${STAGE:?}/deployment"
    if test -z "$PROGRAM_NAME"; then
        ansible-playbook --forks="20" playbook.yml
    else
        ansible-playbook --forks="20" --limit="$PROGRAM_NAME" --tags="redeploy" playbook.yml
    fi
    test "$(curl -sSL ${PING_URL:?})" = "pong" || error "API gateway didn't pong to ping"
)

re() (
    set -x
    down-app "$@"
    up-app "$@"
    deploy ""
)

all() {
    up-vpn --auto-approve
    up-app --auto-approve
    ssh-key-update
    deploy
}

# MARK: VPN

vpn-connect() {
    REGION_SLUG="$1" && shift
    sudo -v
    sudo openvpn "${STAGE:?}/artifacts/vpn/dth-do-${REGION_SLUG:?}-provisioner.ovpn"
    # sleep 1 && sudo killall mDNSResponder{,Helper}
    sudo -k
}

# MARK: Secrets

recreate-certificate-authority() (
    PS4="\n> "
    set -x -T -v -e -E
    mkdir -p "${STAGE:?}/secrets"
    cd "${STAGE:?}/secrets"
    test -d pki && rm -rfv pki
    easyrsa init-pki soft
    easyrsa --batch --req-cn="Logbook Stage Environment CA" build-ca nopass
    security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db "${STAGE:?}/secrets/pki/ca.crt" # macos keychain
)

rotate-cryptographic-keys() (
    PS4="\n> "
    set -x -T -v -e -E
    mkdir -p "${STAGE:?}/secrets"
    cd "${STAGE:?}/secrets"

    ssh-keys() {
        mkdir -p "${STAGE:?}/secrets/ssh"
        test -f "${STAGE:?}/secrets/ssh/application-server" && rm -rfv "${STAGE:?}/secrets/ssh/application-server"
        ssh-keygen -a 1000 -b 4096 -C "application-server" -o -t rsa -f "${STAGE:?}/secrets/ssh/application-server" -N ''
        cp "${STAGE:?}/secrets/ssh/application-server" "${STAGE:?}/image/application/upload/map/home/.ssh/application-server"
        cp "${STAGE:?}/secrets/ssh/application-server.pub" "${STAGE:?}/image/database/upload/map/home/.ssh/authorized_keys"
    }
    ssh-keys

    pki() {
        rotate-server-cert() {
            COMMON_NAME="${1:?}"
            # https://github.com/OpenVPN/easy-rsa/blob/master/doc/EasyRSA-Renew-and-Revoke.md
            if test -f "${PKI:?}/issued/${COMMON_NAME:?}.crt"; then
                if test -f "${PKI:?}/expired/${COMMON_NAME:?}.crt"; then
                    easyrsa --batch revoke-expired "${COMMON_NAME:?}" unspecified
                fi
                easyrsa --batch expire "${COMMON_NAME:?}"
                easyrsa --batch sign-req server "${COMMON_NAME:?}"
            else
                easyrsa --subject-alt-name="DNS:${COMMON_NAME:?}" --batch build-server-full "${COMMON_NAME:?}" nopass
            fi
        }
        local PKI
        PKI="${STAGE:?}/secrets/pki"

        rotate-server-cert "stage.logbook.balaasad.com"

        cp "${PKI:?}/issued/stage.logbook.balaasad.com.crt" \
            "${STAGE:?}/image/gateway/upload/map/etc/ssl/certs/stage.logbook.balaasad.com.crt"
        cp "${PKI:?}/private/stage.logbook.balaasad.com.key" \
            "${STAGE:?}/image/gateway/upload/map/etc/ssl/private/stage.logbook.balaasad.com.key"
    }
    pki
)

# MARK: Digitalocean

do-images() {
    ALL="$(doctl compute snapshot list --format CreatedAt,ID,Name)"
    ALL_P="$(echo "$ALL" | tail -n +2 | gsed -E -e 's/[\ \t]+/ /g' -e 's/_/ /g')"
    KINDS="$(echo "$ALL_P" | cut -d ' ' -f 4 | sort | uniq)"
    UPTODATE_IDs="$(echo "$KINDS" | while read KIND; do echo "$ALL_P" | grep "$KIND" | tail -n 1 | cut -d ' ' -f 2; done)"
    echo "$ALL" | while read IMAGE; do if echo "$IMAGE" | grep "$UPTODATE_IDs" >/dev/null 2>&1; then echo "$IMAGE" | green; else echo "$IMAGE"; fi; done
}

do-up-to-date-images() {
    local ALL
    ALL="$(doctl compute snapshot list --format Name,ID --no-header | grep build | grep -v base | sort -r | gsed -E 's/[ ]+/ /g')"
    echo "$ALL" | gsed -E 's/(build_([^_]+)_([0-9_]{17}).*)/\1 \2/g' | sort -r | uniq -f 2 | cut -d ' ' -f 1-2 | gsed -E 's/^([^\ ]+) (.*)/\2 \1/g'
}

do-clean-outdated-images() {
    set -E
    local ALL
    local TOKEEP
    local OUTDATED
    local OUTDATED_NAMEs
    local OUTDATED_IDs
    ALL="$(doctl compute snapshot list --format Name,ID --no-header | grep build | grep -v base | sort -r | gsed -E 's/[ ]+/ /g')"
    TOKEEP="$(echo "$ALL" | gsed -E 's/(build_([^_]+)_([0-9_]{17}).*)/\1 \2/g' | sort -r | uniq -f 2 | cut -d ' ' -f 1)"
    OUTDATED="$(echo "$ALL" | grep -v "$TOKEEP")"
    OUTDATED_NAMEs="$(echo "$OUTDATED" | cut -d ' ' -f 1)"
    OUTDATED_IDs="$(echo "$OUTDATED" | cut -d ' ' -f 2)"
    test -z "$OUTDATED_IDs" && echo "nothing to clean" && return 0
    echo "keep:" | green
    echo "$TOKEEP" | green
    echo "delete:" | red
    echo "$OUTDATED_NAMEs" | red
    confirm "$ALL" "$OUTDATED_IDs" && echo "$OUTDATED_IDs" | xargs --verbose -n 1 doctl compute image delete -f
}

image-tree() {
    find image -name 'dr.yml' | while read FOLDER; do
        CHILD="$(basename "$(dirname "$FOLDER")")"
        PARENT="$(basename "$(cat "$FOLDER" | yq -r '.depends_on.folder')")"
        echo "$PARENT -> $CHILD"
    done | sort
}
