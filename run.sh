#!/usr/bin/env bash

PID_LIST=()

cd src/upf/build
sudo -E ./bin/free5gc-upfd &
PID_LIST+=($!)

sleep 1

cd ../../..

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

export GIN_MODE=release

for NF in ${NF_LIST}; do
    ./bin/${NF} &
    PID_LIST+=($!)
done

sudo ./bin/n3iwf &
PID_LIST+=($!)

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
