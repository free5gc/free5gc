#!/usr/bin/env bash

# Check OS
if [ -f /etc/os-release ]; then
    # freedesktop.org and systemd
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
else
    # Fall back to uname, e.g. "Linux <version>", also works for BSD, etc.
    OS=$(uname -s)
    VER=$(uname -r)
    echo "This Linux version is too old: $OS:$VER, we don't support!"
    exit 1
fi

sudo -v
if [ $? == 1 ]
then
    echo "Without root permission, you cannot run the test due to our test is using namespace"
    exit 1
fi

UPF_NUM=6

while getopts 'om:' OPT;
do
    case $OPT in
        o) DUMP_NS=True;;
        m)
            SCALE=$OPTARG
            if ! [[ $SCALE =~ ^[0-9]+$ ]]; then
                echo "-m should be number"
                exit 1
            fi
            if [[ $(( $SCALE )) -ge ${UPF_NUM} && $(( $SCALE )) -le 99 ]]; then
                UPF_NUM=$(( $SCALE ))
            else
                echo "UL-CL UPF must larger then ${UPF_NUM} and less then 100"
                exit 1
            fi
            ;;
    esac
done
shift $(($OPTIND - 1))

TEST_POOL="TestULCLAndMultiUPF"
if [[ ! "$1" =~ $TEST_POOL ]]
then
    echo "Usage: $0 [ ${TEST_POOL//|/ | } ]"
    exit 1
fi

GOPATH=$HOME/go
if [ $OS == "Ubuntu" ]; then
    GOROOT=/usr/local/go
elif [ $OS == "Fedora" ]; then
    GOROOT=/usr/lib/golang
fi
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

UPFNS="UPFns"
CONF_DIR=$(cd config && pwd)

export GIN_MODE=release

# Setup bridge
sudo ip link add veth0 type veth peer name br-veth0
sudo ip link set veth0 up
# sudo ip addr add 10.60.0.1 dev lo
sudo ip addr add 10.200.200.1/24 dev veth0
sudo ip addr add 10.200.200.2/24 dev veth0

sudo ip link add free5gc-br type bridge
sudo ip link set free5gc-br up
sudo ip link set br-veth0 up
sudo ip link set br-veth0 master free5gc-br

sudo iptables -I FORWARD 1 -j ACCEPT

# Setup network namespace
for i in $(seq -f "%02g" 1 $UPF_NUM); do
    sudo ip netns add "${UPFNS}${i}"

    sudo ip link add "veth${i}" type veth peer name "br-veth${i}"
    sudo ip link set "br-veth${i}" up
    sudo ip link set "veth${i}" netns "${UPFNS}${i}"

    sudo ip netns exec "${UPFNS}${i}" ip link set lo up
    sudo ip netns exec "${UPFNS}${i}" ip link set "veth${i}" up
#    sudo ip netns exec "${UPFNS}${i}" ip addr add "10.60.0.1${i}" dev lo
    sudo ip netns exec "${UPFNS}${i}" ip addr add "10.200.200.1${i}/24" dev "veth${i}"

    sudo ip link set "br-veth${i}" master free5gc-br

    if [ ${DUMP_NS} ]; then
        sudo ip netns exec "${UPFNS}${i}" tcpdump -i any -w "${UPFNS}${i}.pcap" &
        sleep 1
        TCPDUMP_PID_[${i}]=$(sudo ip netns pids "${UPFNS}${i}")
    fi

    cd NFs/upf/build && sudo -E ip netns exec "${UPFNS}${i}" ./bin/free5gc-upfd -c "${CONF_DIR}/multiUPF/upfcfg${i}.yaml" &
    sleep 1
done

NF_LIST="nrf amf udr pcf udm nssf ausf"
F5GC_DIR="$(cd "$( dirname "$0" )" && pwd -P)"
for NF in ${NF_LIST}; do
    $F5GC_DIR/bin/${NF} -c "${CONF_DIR}/${NF}cfg.yaml"&
    PID_LIST+=($!)
    sleep 0.1
done

$F5GC_DIR/bin/smf -c "${CONF_DIR}/multiUPF/smfcfg.ulcl.yaml" -u "${CONF_DIR}/multiUPF/uerouting.yaml"&
PID_LIST+=($!)

cd test
$GOROOT/bin/go test -v -vet=off -run $1  -args noinit

for ((idx=${#PID_LIST[@]}-1;idx>=0;idx--)); do
    sudo kill -SIGINT ${PID_LIST[$idx]}
    sleep 0.1
done

sleep 3
sudo killall -15 free5gc-upfd
sleep 1

cd ..
mkdir -p testkeylog
for KEYLOG in $(ls *sslkey.log); do
     mv $KEYLOG testkeylog
done

# sudo ip addr del 10.60.0.1/32 dev lo
sudo ip link del veth0
sudo ip link del free5gc-br

sudo iptables -D FORWARD -j ACCEPT

for i in $(seq -f "%02g" 1 $UPF_NUM); do
  if [ ${DUMP_NS} ]; then
      sudo ip netns exec "${UPFNS}${i}" kill -SIGINT ${TCPDUMP_PID_[$i]}
  fi

  sudo ip netns del "${UPFNS}${i}"
  sudo ip link del "br-veth${i}"
done
