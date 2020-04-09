#!/bin/bash
source ./config

cd ${LIBGTP5GNL_TOOLS_PATH}

sudo killall -9 gtp5g-link
sudo ./gtp5g-link add gtp5gtest &
sleep 0.2

sudo ./gtp5g-tunnel add far gtp5gtest 1 --action 2
sudo ./gtp5g-tunnel add far gtp5gtest 2 --action 2 --hdr-creation 1 89 ${IUPF_IP} 2152
sudo ./gtp5g-tunnel add pdr gtp5gtest 1 --pcd 1 --hdr-rm 0 --ue-ipv4 ${UE_IP} --f-teid 88 ${AUPF_IP} --far-id 1
sudo ./gtp5g-tunnel add pdr gtp5gtest 2 --pcd 2 --ue-ipv4 ${UE_IP} --far-id 2

sudo ip r add ${UE_CIDR} dev gtp5gtest
