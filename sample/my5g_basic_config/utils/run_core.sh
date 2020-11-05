#!/usr/bin/env bash
#
# Script not ready yet
#

CFG_DIR="$(pwd)/.."
PID_LIST=()

cd ../../../src/upf/build
sudo -E ip netns exec UPFns ./bin/free5gc-upfd " -free5gccfg $CFG_DIR/my5Gcore.conf" &
PID_LIST+=($!)

sleep 1

cd ../../..

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

export GIN_MODE=release

for NF in ${NF_LIST}; do
    ./bin/${NF} -free5gccfg "$CFG_DIR/my5Gcore.conf" -${NF}cfg "$CFG_DIR/${NF}cfg.conf" &
    PID_LIST+=($!)
done

sudo ./bin/n3iwf -free5gccfg "$CFG_DIR/my5Gcore.conf" -n3iwfcfg "$CFG_DIR/n3iwfcfg.conf" &
SUDO_N3IWF_PID=$!
sleep 1
N3IWF_PID=$(pgrep -P $SUDO_N3IWF_PID)
PID_LIST+=($SUDO_N3IWF_PID $N3IWF_PID)

function terminate()
{
    # kill amf first
    while $(sudo kill -SIGINT ${PID_LIST[2]} 2>/dev/null); do
        sleep 2
    done

    for ((idx=${#PID_LIST[@]}-1;idx>=0;idx--)); do
        sudo kill -SIGKILL ${PID_LIST[$idx]}
    done
}

trap terminate SIGINT
wait ${PID_LIST}
