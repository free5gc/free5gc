#!/bin/bash

NF_LIST="nrf amf smf udr pcf udm nssf ausf n3iwf free5gc-upfd"

for NF in ${NF_LIST}; do
    sudo killall -9 ${NF}
done

sudo killall tcpdump
sudo ip link del upfgtp
sudo ip link del ipsec0
sudo ip link del xfrmi-default
sudo rm /dev/mqueue/*
sudo rm -f /tmp/free5gc_unix_sock
mongo --eval "db.NfProfile.drop()" free5gc

