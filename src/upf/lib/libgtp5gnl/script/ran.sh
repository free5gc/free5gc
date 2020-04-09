#!/bin/bash
source ./config

sudo ip addr add ${UE_IP} dev lo >> /dev/null 2>&1

cd ${LIBGTP5GNL_TOOLS_PATH}

sudo killall -9 gtp5g-link
sudo ./gtp5g-link add gtp5gtest --ran &
sleep 0.2

sudo ./gtp5g-tunnel add far gtp5gtest 1 --action 2
sudo ./gtp5g-tunnel add far gtp5gtest 2 --action 2 --hdr-creation 0 78 ${UPF_IP} 2152
sudo ./gtp5g-tunnel add pdr gtp5gtest 1 --pcd 1 --hdr-rm 0 --ue-ipv4 ${UE_IP} --f-teid 87 ${RAN_IP} --far-id 1
sudo ./gtp5g-tunnel add pdr gtp5gtest 2 --pcd 2 --ue-ipv4 ${UE_IP} --far-id 2

sudo ip r add ${DN_CIDR} dev gtp5gtest
