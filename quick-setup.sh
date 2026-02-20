#!/bin/bash

GO_VERSION="1.25.5"

NETWORK_INTERFACE=""
IP_ADDRESS=""

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
    echo
    echo "Example:"
    echo "  $0 -i eno1"
    echo
    exit 1
}

check_interface_and_get_ip() {
    if ! ip link show "$1" &>/dev/null; then
        log_error "Interface $1 does not exist."
        usage
    fi
    NETWORK_INTERFACE=$1
    IP_ADDRESS=$(ip addr show "$1" | grep -oP '(?<=inet )\S+' | head -n 1 | cut -d/ -f1)
    log_info "Use interface $1 as N2/N3/N6 interface with IP address $IP_ADDRESS."
}

install_golang() {
    log_info "Installing golang..."

    if go version > /dev/null 2>&1; then
        current_version=$(go version | awk '{print $3}' | sed 's/^go//')
        log_info "Go $current_version has already installed."

        if [[ $current_version != $GO_VERSION ]]; then
            upgrade="y"
            log_question "Go $current_version is already installed. Do you want to upgrade to $GO_VERSION? [Y/n] "
            read upgrade
            if [[ $upgrade == "Y" || $upgrade == "y" ]]; then
                sudo rm -rf /usr/local/go
                wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
                sudo tar -C /usr/local -zxvf go${GO_VERSION}.linux-amd64.tar.gz
                rm go${GO_VERSION}.linux-amd64.tar.gz
            fi
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
}

main() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--interface)
                if [[ -z $2 ]]; then
                    usage
                fi
                check_interface_and_get_ip "$2"
                shift 2
                ;;
            *)
                usage
                ;;
        esac
    done
    if [[ -z $NETWORK_INTERFACE ]]; then
        log_warn "Network interface is not set."
    fi
    separate_stars

    install_golang
    separate_stars
}

main "$@"