#!/bin/bash

if [[ $# -ne 2 ]] || ([[ $1 != "up" ]] && [[ $1 != "down" ]]); then
  echo "Usage: $0 [up|down] [internet_iface]"
  exit 1
fi

HOSTNAMES=(
  "10.1.1.2 amf"
  "10.1.1.3 smf"
  "10.1.1.4 ausf"
  "10.1.1.5 nssf"
  "10.1.1.6 pcf"
  "10.1.1.7 udm"
  "10.1.1.8 udr"
  "10.1.1.9 upf"
  "10.1.1.11 db"
  "10.1.1.10 nrf"
)

UENS="UEns"
EXEC_UENS="sudo ip netns exec ${UENS}"

UPFNS="UPFns"
EXEC_UPFNS="sudo ip netns exec ${UPFNS}"
IFACE=$2

if [[ $1 == "up" ]]; then
  echo "Creating network interfaces and namespaces..."
  # create network interfaces and add ip addresses
  # 5gc network (it's not needed but helps to organize/separate the networks)
  sudo ip link add br-5gc type  bridge
  sudo ip addr add 10.1.1.2/24  dev br-5gc
  sudo ip addr add 10.1.1.3/24  dev br-5gc
  sudo ip addr add 10.1.1.4/24  dev br-5gc
  sudo ip addr add 10.1.1.5/24  dev br-5gc
  sudo ip addr add 10.1.1.6/24  dev br-5gc
  sudo ip addr add 10.1.1.7/24  dev br-5gc
  sudo ip addr add 10.1.1.8/24  dev br-5gc
  sudo ip addr add 10.1.1.9/24  dev br-5gc
  sudo ip addr add 10.1.1.10/24 dev br-5gc
  sudo ip addr add 10.1.1.11/24 dev br-5gc
  sudo ip link set br-5gc up

  # Inteface added to handle N2 interface (it's not needed but helps to organize/separate the networks)
  sudo ip link add br-n2 type bridge
  sudo ip addr add 172.16.0.1/24 dev br-n2
  sudo ip addr add 172.16.0.2/24 dev br-n2

  # Setup network namespace for UPF
  sudo ip netns add ${UPFNS}

  sudo ip link add veth0 type veth peer name veth1
  sudo ip link set veth0 up
  sudo ip addr add 60.60.0.1 dev lo
  sudo ip addr add 10.200.200.1/24 dev veth0
  sudo ip addr add 10.200.200.2/24 dev veth0

  sudo ip link set veth1 netns ${UPFNS}
  ${EXEC_UPFNS} ip link set lo up
  ${EXEC_UPFNS} ip link set veth1 up
  ${EXEC_UPFNS} ip addr add 60.60.0.101 dev lo
  ${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
  ${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1

  #sudo ip link set dev upfgtp mtu 1500
  #${EXEC_UPFNS} ip link set dev upfgtp mtu 1500

  sudo ip netns add ${UENS}
  sudo ip link add veth2 type veth peer name veth3
  sudo ip addr add 192.168.127.1/24 dev veth2
  sudo ip link set veth2 up

  sudo ip link set veth3 netns ${UENS}
  ${EXEC_UENS} ip addr add 192.168.127.2/24 dev veth3
  ${EXEC_UENS} ip link set lo up
  ${EXEC_UENS} ip link set veth3 up
  ${EXEC_UENS} ip link add ipsec0 type vti local 192.168.127.2 remote 192.168.127.1 key 5
  ${EXEC_UENS} ip link set ipsec0 up

  sudo ip link add name ipsec0 type vti local 192.168.127.1 remote 0.0.0.0 key 5
  sudo ip addr add 10.0.0.1/24 dev ipsec0
  sudo ip link set ipsec0 up

  sudo ip link add veth4 type veth peer name veth5
  sudo ip addr add 10.1.2.1/24 dev veth4
  sudo ip link set veth4 up

  sudo ip link set veth5 netns ${UPFNS}
  ${EXEC_UPFNS} ip addr add 10.1.2.2/24 dev veth5
  ${EXEC_UPFNS} ip link set veth5 up
  ${EXEC_UPFNS} ip route add default via 10.1.2.1

  ${EXEC_UPFNS} iptables -t nat -A POSTROUTING -o veth5 -j MASQUERADE
  sudo iptables -t nat -A POSTROUTING -s 10.1.2.2/24 -o ${IFACE} -j MASQUERADE
  sudo iptables -A FORWARD -i ${IFACE} -o veth4 -j ACCEPT
  sudo iptables -A FORWARD -o ${IFACE} -i veth4 -j ACCEPT

  sudo sysctl -w net.ipv4.ip_forward=1 > /dev/null
  echo "Network interfaces and namespaces created."

elif [[ $1 == "down" ]]; then
  echo "Removing network interfaces and namespaces.."
  sudo ip link set br-5gc down
  sudo ip link delete br-5gc
  sudo ip link delete br-n2
  sudo ip xfrm policy flush
  sudo ip xfrm state flush
  sudo ip link del veth2
  sudo ip link del veth4
  sudo ip link del ipsec0
  sudo ip link del veth0
  ${EXEC_UENS} ip link del ipsec0
  sudo ip netns del ${UENS}
  sudo ip netns del ${UPFNS}

  sudo rm /dev/mqueue/*
  for host in "${HOSTNAMES[@]}"; do
    sudo sed -i "/$host/d" /etc/hosts
  done
  echo "Network interfaces and namespaces removed."
fi