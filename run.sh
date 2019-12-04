#!/usr/bin/env bash

cd src/upf/build
sudo ./bin/free5gc-upfd &
PID_LIST=$!

cd ../../..

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

for NF in ${NF_LIST}; do 
    ./bin/${NF} &
    PID_LIST="${PID_LIST} $!"
done

trap "sudo kill -SIGKILL ${PID_LIST}" SIGINT
wait ${PID_LIST}