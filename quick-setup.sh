#!/bin/bash

N6_INTERFACE=""

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
    echo -e "${COLOR_YELLOW}[*]${COLOR_RESET} $1"
}

log_error() {
    echo -e "${COLOR_RED}[-]${COLOR_RESET} $1"
}

usage() {
    echo "usage: $0 [options]"
    echo "  -n6, --n6 <interface>   Setup the n6 network interface for upf to forward the packets."

    exit 1
}

main() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n6|--n6)
                if [[ -z $2 ]]; then
                    usage
                fi
                N6_INTERFACE=$2
                log_info "Use ${N6_INTERFACE} as N6 interface."
                shift 2
                ;;
            *)
                usage
                ;;
        esac
    done
}

main "$@"