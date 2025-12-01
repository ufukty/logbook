#!/usr/bin/env bash

set -xe
PS4='\033[31m$0:$LINENO: \033[0m'

test -d .git || (echo "Run from root folder" && exit 1)

test -f ".venv/bin/activate" ||
  python3 -m venv ".venv"
# shellcheck disable=SC1091
source ".venv/bin/activate"

which make ||
  xcode-select --install

which go ||
  (open "https://go.dev/dl" && exit 1)

which stringer ||
  go install "golang.org/x/tools/cmd/stringer@latest"
which gonfique ||
  go install "github.com/ufukty/gonfique@v1.3.1"
which sqlc ||
  go install "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
which govalid ||
  go install "github.com/ufukty/govalid@v0.1.0"
which d2 ||
  go install "oss.terrastruct.com/d2@v0.6.3"
which gohandlers ||
  go install "github.com/ufukty/gohandlers/cmd/gohandlers@latest"
which ovpn-auth ||
  go install "github.com/ufukty/ovpn-auth/cmd/ovpn-auth@v1.1.1"
which shfmt ||
  go install mvdan.cc/sh/v3/cmd/shfmt@latest

(bash --version | grep "^GNU bash, version 5") ||
  brew install "bash"
which gsed ||
  brew install gnu-sed
which psql ||
  (brew install "postgresql@16" && brew services start "postgresql@16")
which doctl ||
  brew install "doctl"
(which terraform && which packer) ||
  brew tap "hashicorp/tap"
which terraform ||
  brew install "hashicorp/tap/terraform"
which packer ||
  brew install "hashicorp/tap/packer"
which openvpn ||
  brew install "openvpn"
which easyrsa || # maintain PKI
  (open "https://github.com/OpenVPN/easy-rsa" && exit 1)
which jq || # platform
  brew install jq
which unbuffer || # run.sh (to trick Chi logger to print colors)
  brew install expect
which shellcheck ||
  brew install shellcheck
which envsubst ||
  brew install gettext

(which ansible && which qr) ||
  pip install -r "dev/requirements.txt"

which argon2 ||
  (open "https://github.com/P-H-C/phc-winner-argon2/releases/tag/20190702" && exit 1)

which npm ||
  (open "https://nodejs.org/en/download" && exit 1)
which mmdc || # docs
  npm install -g "@mermaid-js/mermaid-cli"

test -f ~/.bash_include/autosource.sh ||
  (mkdir -p ~/.bash_include && cp dev/data/autosource.sh ~/.bash_include/autosource.sh)

# shellcheck disable=2016
grep ". ~/.bash_include" ~/.bash_profile >/dev/null ||
  echo 'for f in ~/.bash_include/*.sh; do . $f; done' >>~/.bash_profile

# shellcheck disable=SC2016
test -d platform/stage/provision/application/.terraform/providers ||
  (cd platform/stage/provision/application && terraform init)

# shellcheck disable=SC2016
test -d platform/stage/provision/vpc/.terraform/providers ||
  (cd platform/stage/provision/vpc && terraform init)

# shellcheck disable=SC2016
test -d platform/stage/provision/vpn/.terraform/providers ||
  (cd platform/stage/provision/vpn && terraform init)
