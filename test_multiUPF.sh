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
            if [[ $(( SCALE )) -ge ${UPF_NUM} && $(( SCALE )) -le 99 ]]; then
                UPF_NUM=$(( SCALE ))
            else
                echo "UL-CL UPF must larger then ${UPF_NUM} and less then 100"
                exit 1
            fi
            ;;
        *) ;;
    esac
done
shift $((OPTIND - 1))

TEST_POOL="TestULCLAndMultiUPF"
if [[ ! "$1" =~ $TEST_POOL ]]
then
    echo "Usage: $0 [ ${TEST_POOL//|/ | } ]"
    exit 1
fi

GOPATH=$HOME/go
if [[ $OS == "Ubuntu" ]]; then
    GOROOT=/usr/local/go
elif [[ $OS == "Fedora" ]]; then
    GOROOT=/usr/lib/golang
fi
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

UPFNS="UPFns"
CONF_DIR=$(cd config && pwd)

export GIN_MODE=release

# Setup bridge
sudo ip link add veth0 type veth peer name br-veth0
sudo ip link set veth0 up
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
    sudo ip netns exec "${UPFNS}${i}" ip addr add "10.200.200.1${i}/24" dev "veth${i}"

    sudo ip link set "br-veth${i}" master free5gc-br

    if [ ${DUMP_NS} ]; then
        sudo ip netns exec "${UPFNS}${i}" tcpdump -i any -w "${UPFNS}${i}.pcap" &
    fi

    sudo -E ip netns exec "${UPFNS}${i}" ./bin/upf -c "${CONF_DIR}/multiUPF/upfcfg${i}.yaml" &
    sleep 1
done

if [ ${DUMP_NS} ]; then
    sudo tcpdump -i any 'sctp port 38412 || tcp port 8000 || udp port 8805' -w 'control_plane.pcap' &
fi

NF_LIST="nrf udr udm ausf nssf amf pcf"
F5GC_DIR="$(cd "$( dirname "$0" )" && pwd -P)"
for NF in ${NF_LIST}; do
    "${F5GC_DIR}/bin/${NF}" -c "${CONF_DIR}/${NF}cfg.yaml"&
    sleep 0.1
done
"${F5GC_DIR}/bin/smf" -c "${CONF_DIR}/multiUPF/smfcfg.ulcl.yaml" -u "${CONF_DIR}/multiUPF/uerouting.yaml"&
NF_LIST+=" smf"

# git clone go-gtp5gnl and build gogtp5g-tunnel
# gogtp5g-tunnel will be used to show PDR/FAR during testing
./make_gtp5gtunnel.sh


# ensure to kill remaining processes and to remove addresses for the test
function cleanup {
    local procs

    pkill -e -SIGTERM amf
    procs="$(echo "${NF_LIST}" | sed -E -e 's/ *nrf *//' -e 's/ *amf */ /' -e 's/ +/|/g')"
    pkill -e -SIGTERM "(${procs})"
    sleep 1
    pkill -e nrf

    sudo killall -SIGTERM upf
    sleep 2

    if [ ${DUMP_NS} ]; then
        sudo killall tcpdump
    fi
    for i in $(seq -f "%02g" 1 ${UPF_NUM}); do
      sudo ip link del "br-veth${i}"
      sudo ip netns del "${UPFNS}${i}"
    done

    sudo ip link del veth0
    sudo ip link del free5gc-br
    sudo iptables -D FORWARD -j ACCEPT
    sleep 2
}

trap cleanup EXIT

cd test && $GOROOT/bin/go test -v -vet=off -run "^$1$" -args noinit
