#!/bin/bash

GO_VERSION="1.25.5"
GOLANGCI_LINT_VERSION="2.7.2"

NETWORK_INTERFACE=""
IP_ADDRESS=""
LINT=false
DOCKER=false

GTP5G_PATH="$HOME/gtp5g"

SUCCESS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0

COLOR_RED="\033[31m"
COLOR_GREEN="\033[32m"
COLOR_YELLOW="\033[33m"
COLOR_BLUE="\033[36m"
COLOR_RESET="\033[0m"

log_info() {
    echo -e "${COLOR_BLUE}[.]${COLOR_RESET} $1"
}

log_success() {
    echo -e "${COLOR_GREEN}[+]${COLOR_RESET} $1"
}

log_warn() {
    echo -e "${COLOR_YELLOW}[!]${COLOR_RESET} $1"
}

log_question() {
    echo -e "${COLOR_YELLOW}[?]${COLOR_RESET} $1"
}

log_error() {
    echo -e "${COLOR_RED}[-]${COLOR_RESET} $1"
}

separate_stars() {
    local cols=$(tput cols 2>/dev/null || echo 10)
    printf "%*s\n" "$cols" "" | tr ' ' '*'
}

usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -i, --interface <network interface>   Setup the network interface for N2/N3/N6"
    echo "  -l, --lint                            Enable golangci-lint installation"
    echo "  -d, --docker                          Enable docker installation"
    echo "  -h, --help                            Show this help message"
    echo
    echo "Example:"
    echo "  $0 -i eno1 -l -d"
    echo
}

check_interface_and_get_ip() {
    if ! ip link show "$1" &>/dev/null; then
        log_error "Interface $1 does not exist."
        return 1
    fi
    NETWORK_INTERFACE=$1
    IP_ADDRESS=$(ip addr show "$1" | grep -oP '(?<=inet )\S+' | head -n 1 | cut -d/ -f1)
    log_info "Use interface $1 as N2/N3/N6 interface with IP address $IP_ADDRESS."
}

network_config() {
    log_info "Configuring network..."

    sudo systemctl stop ufw
    sudo systemctl disable ufw

    sudo sysctl -w net.ipv4.ip_forward=1
    sudo iptables -t nat -A POSTROUTING -o ${NETWORK_INTERFACE} -j MASQUERADE
    sudo iptables -I FORWARD 1 -j ACCEPT

    log_success "Network configured with interface ${NETWORK_INTERFACE}"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_golang() {
    log_info "Installing golang..."

    if go version > /dev/null 2>&1; then
        current_version=$(go version | awk '{print $3}' | sed 's/^go//')
        log_info "Go $current_version already installed."

        if [[ $current_version != $GO_VERSION ]]; then
            log_question "Go $current_version is already installed. Do you want to upgrade to $GO_VERSION? [Y/n] "
            read upgrade
            upgrade=${upgrade:-Y}
            if [[ $upgrade == "Y" || $upgrade == "y" ]]; then
                sudo rm -rf /usr/local/go
                wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
                sudo tar -C /usr/local -zxvf go${GO_VERSION}.linux-amd64.tar.gz
                rm go${GO_VERSION}.linux-amd64.tar.gz
            fi

            log_success "Golang upgraded"
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        else
            SKIP_COUNT=$((SKIP_COUNT + 1))
        fi
        return
    fi

    wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
    sudo tar -C /usr/local -zxvf go${GO_VERSION}.linux-amd64.tar.gz
    mkdir -p ~/go/{bin,pkg,src}
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
    echo 'export GO111MODULE=auto' >> ~/.bashrc
    source ~/.bashrc
    rm go${GO_VERSION}.linux-amd64.tar.gz

    log_success "Golang installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_golangci_lint() {
    log_info "Installing golangci-lint..."

    if golangci-lint version > /dev/null 2>&1; then
        log_info "Golangci-lint already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi
    curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin $GOLANGCI_LINT_VERSION

    log_success "Golangci-lint installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_mongodb() {
    log_info "Installing mongoDB..."

    if mongosh --version > /dev/null 2>&1; then
        log_info "MongoDB already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi

    sudo apt-get install gnupg curl
    curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg --dearmor

    ubuntu_version=$(lsb_release -rs)
    case $ubuntu_version in
        20.04)
            echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/8.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-8.2.list
            ;;
        22.04)
            echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/8.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-8.2.list
            ;;
        24.04)
            echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu noble/mongodb-org/8.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-8.2.list
            ;;
        25.04)
            echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu noble/mongodb-org/8.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-8.2.list
            ;;
        *)
            log_error "Unsupported Ubuntu version: $ubuntu_version"
            FAIL_COUNT=$((FAIL_COUNT + 1))
            return 1
            ;;
    esac
    sudo apt-get update
    sudo apt-get install -y mongodb-org
    sudo systemctl enable --now mongod

    log_success "MongoDB installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_gtp5g() {
    log_info "Installing gtp5g..."

    if lsmod | grep -q gtp5g; then
        log_info "GTP5G already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi

    sudo apt -y update
    sudo apt -y install gcc g++ cmake autoconf libtool pkg-config libmnl-dev libyaml-dev

    git clone https://github.com/free5gc/gtp5g.git $GTP5G_PATH
    pushd $GTP5G_PATH
    make
    sudo make install
    popd

    log_success "GTP5G installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_yarn() {
    log_info "Installing Yarn..."

    if yarn --version > /dev/null 2>&1; then
        log_info "Yarn already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - 
    sudo apt update
    sudo apt install -y nodejs
    sudo corepack enable
    echo 'export COREPACK_ENABLE_DOWNLOAD_PROMPT=0' >> ~/.bashrc
    export COREPACK_ENABLE_DOWNLOAD_PROMPT=0

    log_success "Yarn installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

install_docker() {
    log_info "Installing Docker..."

    if docker --version > /dev/null 2>&1; then
        log_info "Docker already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi

    sudo apt update
    sudo apt install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

    sudo apt update
    sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    sudo groupadd docker
    sudo usermod -aG docker $USER

    log_success "Docker installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

submodule_init() {
    log_info "Initializing submodules..."

    git submodule update --init --recursive
    log_success "Submodules initialized"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

make_free5gc() {
    log_info "Making free5GC..."

    make all

    log_success "free5GC made"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

substitute_ip() {
    log_info "Substituting IP addresses in config files..."

    sed -i -e '/ngapIpList:/!b;n;s/- .*/- '"$IP_ADDRESS"'/' config/amfcfg.yaml
    sed -i -e '/endpoints:/!b;n;s/- .*/- '"$IP_ADDRESS"'/' config/smfcfg.yaml
    sed -i -e '/ifList:/!b;n;s/addr: .*/addr: '"$IP_ADDRESS"'/' config/upfcfg.yaml

    log_success "IP addresses substituted"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

print_counts() {
    log_info "Task Summary:"
    log_success "  Success: $SUCCESS_COUNT"
    log_warn "  Skip: $SKIP_COUNT"
    log_error "  Fail: $FAIL_COUNT"
}

main() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--interface)
                if [[ -z $2 ]]; then
                    usage
                    return 1
                fi
                check_interface_and_get_ip "$2" || return 1
                shift 2
                ;;
            -l|--lint)
                LINT=true
                shift
                ;;
            -d|--docker)
                DOCKER=true
                shift
                ;;
            -h|--help)
                usage
                return 0
                ;;
            *)
                usage
                return 1
                ;;
        esac
    done

    if [[ ${NETWORK_INTERFACE} ]]; then
        network_config
        separate_stars

        substitute_ip
        separate_stars
    else
        log_warn "Network interface is not set. Skip network configuration."
        separate_stars
    fi

    install_golang
    separate_stars

    if $LINT; then
        install_golangci_lint
        separate_stars
    fi

    install_mongodb || return 1
    separate_stars

    install_gtp5g
    separate_stars

    install_yarn
    separate_stars

    if $DOCKER; then
        install_docker
        separate_stars
    fi

    submodule_init
    separate_stars

    make_free5gc
    separate_stars

    if [[ ${NETWORK_INTERFACE} ]]; then
        substitute_ip
        separate_stars
    else
        log_warn "Network interface is not set. Skip N2/N3 IP substitution."
        separate_stars
    fi

    print_counts
    separate_stars

    if $DOCKER; then
        log_info "Docker is installed. Please log out and log in again to use Docker without sudo."
    fi
}

main "$@"