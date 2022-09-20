#!/bin/bash
SCRIPT_DIR="$(cd "$( dirname "$0" )" && pwd -P)"
CMD_FILE="gtp5g-tunnel"

[[ -d "${SCRIPT_DIR}/go-gtp5gnl" ]] && rm -fr "${SCRIPT_DIR}/go-gtp5gnl"
[[ -f "${CMD_FILE}" ]] && rm -f "${CMD_FILE}"

git clone https://github.com/free5gc/go-gtp5gnl.git "${SCRIPT_DIR}/go-gtp5gnl"

mkdir "${SCRIPT_DIR}/go-gtp5gnl/bin"
cd "${SCRIPT_DIR}/go-gtp5gnl/cmd/gogtp5g-tunnel" &&  go build -o "${SCRIPT_DIR}/go-gtp5gnl/bin/${CMD_FILE}" .
