#!/usr/local/bin/bash

# MARK: SSH overload

alias ssh="ssh -F $WORKSPACE/platform/stage/artifacts/ssh.conf"
alias scp="scp -F $WORKSPACE/platform/stage/artifacts/ssh.conf"
_check_env_vars() {
    : "${DIGITALOCEAN_ACCESS_TOKEN:?}"
    : "${TF_VAR_DIGITALOCEAN_TOKEN:?}"
    : "${TF_VAR_OVPN_AUTH_USERNAME:?}"
    : "${TF_VAR_OVPN_AUTH_HASH:?}"
    : "${TF_VAR_OVPN_AUTH_TOTP:?}"
}
_check_env_vars

_ssh_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD - 1]}"
    opts=$(grep '^Host' $WORKSPACE/platform/stage/artifacts/ssh.conf 2>/dev/null | grep -v '[?*]' | cut -d ' ' -f 2-)

    COMPREPLY=($(compgen -W "$opts" -- ${cur}))
    return 0
}
complete -F _ssh_completion ssh

# MARK: Utilities

PING_URL="stage.logbook.balaasad.com:8080/api/v1.0.0/ping"

artifacts-update-ssh-aggregation() {
    cd "$WORKSPACE/platform/stage/artifacts"
    cat ssh.conf.d/* >ssh.conf
}

ssh-key-update() {
    touch "${WORKSPACE:?}/platform/stage/artifacts/deployment/service_discovery.json"
    ADDRESSES="$(cat "${WORKSPACE:?}/platform/stage/artifacts/deployment/service_discovery.json" | jq -r '.digitalocean.fra1.services[] | .[] | .ipv4_address_private')"
    echo "$ADDRESSES" | while read ADDRESS; do
        ssh-keygen -R "$ADDRESS" >/dev/null 2>&1
        # ssh-keyscan "$ADDRESS" >>~/.ssh/known_hosts 2>/dev/null
    done
}

update-dns-records() (
    TEMP_FILE="$(mktemp)"
    local GATEWAY_IP
    GATEWAY_IP="$(cat "${WORKSPACE:?}/platform/stage/artifacts/deployment/service_discovery.json" | jq -r '.digitalocean.fra1.services["api-gateway"][0].ipv4_address')"
    ssh -t -F "$WORKSPACE/platform/stage/artifacts/ssh.conf" \
        fra1-vpn "sudo bash -c \"sed \\\"s;{{GATEWAY_IP}};${GATEWAY_IP:?};g\\\" /etc/unbound/unbound.conf.tmpl.d/custom.conf > /etc/unbound/unbound.conf.d/custom.conf && systemctl restart unbound && echo DONE.\""
    test "${OSTYPE:0:6}" = darwin && sudo killall mDNSResponder{,Helper}
)

# MARK: Provision

vpc-up() (
    cd "${WORKSPACE:?}/platform/stage/provisioning/vpc"
    terraform apply "$@" --var-file="${WORKSPACE:?}/platform/stage/provisioning/vars.tfvars"
)

vpn-up() (
    cd "${WORKSPACE:?}/platform/stage/provisioning/vpn"
    terraform apply "$@" --var-file="${WORKSPACE:?}/platform/stage/provisioning/vars.tfvars"
    artifacts-update-ssh-aggregation
    vpn-totp
    read -p "$(note "Connect vpn in separate tab [Enter]")"
    ssh-key-update
    update-dns-records
)

vpn-down() (
    cd "${WORKSPACE:?}/platform/stage/provisioning/vpn"
    terraform destroy "$@" --var-file="${WORKSPACE:?}/platform/stage/provisioning/vars.tfvars"
    artifacts-update-ssh-aggregation
)

app-up() (
    cd "${WORKSPACE:?}/platform/stage/provisioning/application"
    terraform apply "$@" --var-file="${WORKSPACE:?}/platform/stage/provisioning/vars.tfvars"
    artifacts-update-ssh-aggregation
    ssh-key-update
    update-dns-records
)

app-down() (
    cd "${WORKSPACE:?}/platform/stage/provisioning/application"
    terraform destroy "$@" --var-file="${WORKSPACE:?}/platform/stage/provisioning/vars.tfvars"
    artifacts-update-ssh-aggregation
)

# MARK: Deployment

deploy() (
    export PROGRAM_NAME="$1" && shift
    cd "${WORKSPACE:?}/platform/stage/deployment"
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

vpn-totp() (
    cd "${WORKSPACE:?}/platform/stage"
    if test -n "$(find artifacts/vpn -name '*totp*')"; then
        cd "${WORKSPACE:?}/platform/stage/artifacts/vpn"
        find . -name '*totp*' | xargs -n 1 cat | qr
        find . -name '*totp*' -delete
    fi
)

vpn-connect() {
    REGION_SLUG="$1" && shift
    sudo -v
    sudo openvpn "${WORKSPACE:?}/platform/stage/artifacts/vpn/dth-do-${REGION_SLUG:?}-provisioner.ovpn"
    # sleep 1 && sudo killall mDNSResponder{,Helper}
    sudo -k
}

# MARK: Secrets

generate-ca() (
    cd "${WORKSPACE:?}/platform/stage/secrets"
    easyrsa init-pki soft
    easyrsa --batch --req-cn="logbook-CA" build-ca nopass
)

generate-keys() (
    ssh-app-db() (
        set -e
        set -x
        mkdir -p "image/ssh-app-db" && cd "image/ssh-app-db"
        ssh-keygen -a 1000 -b 4096 -C "ssh-app-db" -o -t rsa -f app-db -N '' >/dev/null
    )
    tls-application() (
        set -e
        set -x
        easyrsa --batch build-server-full logbook-application nopass
    )
    tls-non-specific() (
        set -e
        set -x
        easyrsa --batch build-server-full logbook-non-specific nopass
    )
    mkdir -p "$WORKSPACE/platform/stage/secrets"
    cd "$WORKSPACE/platform/stage/secrets"
    ssh-app-db
    tls-application
    tls-non-specific
)

# MARK: Digitalocean

image-clean() {
    function print_outdated_images() {
        set -e
        set -x
        LATEST_IMAGES="$(
            echo "$1" |
                sed -E 's/packer-(.*)-([0-9]{2}-[0-9]{2}-[0-9]{2})-T-([0-9]{2}-[0-9]{2}-[0-9]{2})-UTC[^0-9]*([0-9]*)/\2-\3 \4 \1/g' |
                uniq -f 2 | cut -d ' ' -f 2
        )"
        OUTDATED_IMAGES="$(echo "$1" | grep -v "$LATEST_IMAGES" | sed -E "s/(.*)\ +(.*)/\2/g")"
        echo "$OUTDATED_IMAGES"
    }
    function pretty_print() {
        set -e
        set -x
        OUTDATED_LINES="$(echo "$1" | grep "$2")"
        echo -e "Records highlighted with red will be deleted:"
        echo "$1" | grep -E --color "$OUTDATED_LINES"
    }

    echo "Downloading data from API..."
    IMAGES="$(doctl compute image list-user --format='Name,ID' --no-header | grep -v packer-base | sort -r)"
    OUTDATED_IMAGES="$(print_outdated_images "$IMAGES")"
    test -z "$OUTDATED_IMAGES" && echo "All up-to-date." && exit 0
    pretty_print "$IMAGES" "$OUTDATED_IMAGES"
    while true; do
        read -p "Do you want to delete? Options: (D)elete, (A)bort: " -n 1 USER_INPUT
        case $USER_INPUT in
        [Dd]*)
            echo "Delete."
            doctl compute image delete -f "$OUTDATED_IMAGES"
            exit
            ;;
        [Aa]*)
            echo "Aborted."
            exit
            ;;
        *) echo "Invalid entry. Try again." ;;
        esac
    done
}

image-tree() {
    find image -name 'dr.yml' | while read FOLDER; do
        CHILD="$(basename "$(dirname "$FOLDER")")"
        PARENT="$(basename "$(cat "$FOLDER" | yq -r '.depends_on.folder')")"
        echo "$PARENT -> $CHILD"
    done | sort
}
