#!/usr/local/bin/bash

set -e # exit on error

test -f "$HOME/venv/bin/activate" ||
  python3 -m venv "$HOME/venv"
source "$HOME/venv/bin/activate"

test -f ".venv/bin/activate" ||
  python3 -m venv ".venv"
source ".venv/bin/activate"

which make ||
  xcode-select --install

which go ||
  (open "https://go.dev/dl" && exit 1)

which stringer || # backend
  go install "golang.org/x/tools/cmd/stringer@latest"
which gonfique || # backend
  go install "github.com/ufukty/gonfique@v1.3.1"
which sqlc || # backend
  go install "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
which govalid || # backend
  go install "github.com/ufukty/govalid@v0.1.0"
which d2 || # docs
  go install "oss.terrastruct.com/d2@v0.6.3"

which gohandlers ||
  (echo "install gohandlers" && exit 1)

test -f "/usr/local/bin/bash" ||
  brew install "bash"
which gsed ||
  brew install coreutils
which psql ||
  (brew install "postgresql@15" && brew services start "postgresql@15")
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
which magick || # docs (dark versions of schemas)
  brew install imagemagick

(which ansible && which qr) ||
  pip install -r "$WORKSPACE/dependencies/requirements.txt"

which argon2 ||
  (open "https://github.com/P-H-C/phc-winner-argon2/releases/tag/20190702" && exit 1)

which npm ||
  (open "https://nodejs.org/en/download" && exit 1)
which mmdc || # docs
  npm install -g "@mermaid-js/mermaid-cli"
