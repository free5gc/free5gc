#!/bin/bash

if [[ $# -ne 1 ]] || ([[ $1 != "up" ]] && [[ $1 != "down" ]]); then
  echo "Usage: $0 [up|down]"
  exit 1
fi

FREE5GC_DIR="$(pwd)/../"

# TODO: Read IPs from config files
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

KERNEL=$(uname -r)
KERNEL_VERSION=${KERNEL:0:1}

UENS="UEns"
EXEC_UENS="sudo ip netns exec ${UENS}"

UPFNS="UPFns"
EXEC_UPFNS="sudo ip netns exec ${UPFNS}"
IFACE="enp6s0"

if [[ $1 == "up" ]]; then
  echo "Setting up development environment (dependencies, network namespaces and interfaces)"

  ######################################################
  # INSTALL DEPENDENCIES
  ######################################################

  # install upf
#  echo "Building UPF..."
#  mkdir -p "$FREE5GC_DIR/src/upf/build"
#  cd "$FREE5GC_DIR/src/upf/build"
#  cmake ..
#  make -j`nproc`
#  echo "UPF built"

  # dependencies required
  sudo apt-get install linux-headers-$(uname -r)

  # clone gtp5g repo if it does not exists
  echo "Checking if gtp5g is installed..."
  if [[ ! -d "$FREE5GC_DIR/src/upf/build/gtp5g" ]]; then
    echo "gtp5g not installed"
    GTP5G_REPO_OWNER="jplobianco" #bjoern-r"

    if [[ $KERNEL_VERSION -ge 5 ]] ; then
      GTP5G_REPO_OWNER="PrinzOwO"
    fi

    echo "Cloning gtp5g repo"
    git clone https://github.com/$GTP5G_REPO_OWNER/gtp5g.git "$FREE5GC_DIR/src/upf/build/gtp5g"
  else
    echo "gtp5g is already installed"
  fi

  echo "Checking if gtp5g module is loaded in kernel..."
  lsmod | grep gtp5g >/dev/null
  if [ $? == 1 ]; then
    # Load gtp5g modules on the kernel
    echo "gt5g kernel module not loaded"
    echo "Loading gtp5g module into the kernel..."
    cd "$FREE5GC_DIR/src/upf/build/gtp5g"
    make
    sudo make install
    echo "gtp5g kernel module loaded"
  else
    echo "gtp5g kernel module is already loaded"
  fi

  ######################################################
  # SETTING UP NETWORK INTERFACES AND NAMESPACES
  ######################################################

  # adding hostnames to /etc/hosts
#  echo "Adding hostnames to /etc/hosts..."
#  for host in "${HOSTNAMES[@]}"; do
#    grep "$host" /etc/hosts > /dev/null
#    if [[ $? -gt 0 ]]; then
#      echo "$host" >> /etc/hosts > /dev/null
#      echo "Hostname entry '$host' added to /etc/hosts"
#    else
#      echo "Hostname entry '$host' is already specified in /etc/hosts"
#    fi
#  done

  # create network namespaces
  echo "Creating network namespaces.."
  ip netns list | grep $UPFNS
  if [[ $? -ne 0 ]]; then
    sudo ip netns add $UPFNS
  fi
  ip netns list | grep $UENS
  if [[ $? -ne 0 ]]; then
    sudo ip netns add $UENS
  fi
  echo "Network namespaces created"

  # create network interfaces
  echo "Creating network interfaces.."
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

  #export GIN_MODE=release

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

  #sudo ip link set dev upfgtp0 mtu 1500
  #${EXEC_UPFNS} ip link set dev upfgtp0 mtu 1500

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

  # try to connect with internet
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
  # end try to connect with internet

  echo "Network interfaces created"

  echo "Removing free5gc mongo database..."
  #mongo free5gc --eval 'db.dropDatabase()' > /dev/null
  echo "free5gc mongo database removed"

  echo "Add iptables rules to routing UPF incoming data packets to DN "
  sudo sysctl -w net.ipv4.ip_forward=1 > /dev/null
  # TODO: Check this configuration after ping succeed
  #iptables -t nat -A POSTROUTING -o {DN_Interface_Name} -j MASQUERADE
  echo "iptables rules added"

elif [[ $1 == "down" ]]; then
  echo "Cleaning development environment"

  echo "Removing network namespaces and interfaces.."
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
  echo "Network namespace and interfaces removed"

  echo "Removing free5gc mongo database..."
  #mongo free5gc --eval 'db.dropDatabase()' > /dev/null
  echo "free5gc mongo database removed"

  echo "Cleaning and removing gtp5g..."
  cd $FREE5GC_DIR/src/upf/build/gtp5g
  make > /dev/null
  sudo make clean > /dev/null
  lsmod | grep gtp5g > /dev/null
  if [[ $? == 0 ]]; then
    sudo make uninstall > /dev/null
  fi
  echo "gtp5g cleaned and removed from kernel"

#  echo "Removing hostnames from /etc/hosts..."
#  for host in "${HOSTNAMES[@]}"; do
#    sudo sed -i "/$host/d" /etc/hosts
#  done
#  echo "Hostnames removed from /etc/hosts"

  # TODO: remove message queues
  rm /dev/mqueue/*

  # TODO: remove DN routing rule from iptables

  # TODO: add remove gtp5 devices (use libgtp5gnl tools)
fi
