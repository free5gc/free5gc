#!/bin/bash

PID_LIST=()

cd ../../../src/upf/build
sudo ip netns exec UPFns ./bin/free5gc-upfd &
PID_LIST+=($!)

sleep 5
sudo ip netns exec UPFns ip link set dev upfgtp mtu 1500

wait ${PID_LIST}