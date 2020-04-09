#!/bin/bash
source ./config

echo "############## Stopping ##############"
sudo killall -15 gtp5g-link
sleep 1
 
sudo ip link del gtp5gtest
sudo ip link del veth0
sudo ip netns del ${NS2}
sudo ip addr del ${UE_IP}/32 dev lo
sudo ip addr del ${DN_IP}/32 dev lo

