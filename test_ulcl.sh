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

UPF_NUM=2

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

TEST_POOL="TestRequestTwoPDUSessions"
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

export GIN_MODE=release

# Setup bridge
sudo ip link add veth0 type veth peer name br-veth0
sudo ip link set veth0 up
sudo ip addr add 10.60.0.1 dev lo
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
    sudo ip netns exec "${UPFNS}${i}" ip addr add "10.60.0.1${i}" dev lo
    sudo ip netns exec "${UPFNS}${i}" ip addr add "10.200.200.1${i}/24" dev "veth${i}"

    sudo ip link set "br-veth${i}" master free5gc-br

    if [ ${DUMP_NS} ]; then
        sudo ip netns exec "${UPFNS}${i}" tcpdump -i any -w "${UPFNS}${i}.pcap" &
        sleep 1
        TCPDUMP_PID_[${i}]=$(sudo ip netns pids "${UPFNS}${i}")
    fi

    sed -i -e "s/10.200.200.10./10.200.200.1${i}/g" ./config/upfcfg.testulcl.yaml
    if [ ${i} -eq 02 ]; then
        sed -i -e "s/internet/internet2/g" ./config/upfcfg.testulcl.yaml
    else
        sed -i -e "s/internet2/internet/g" ./config/upfcfg.testulcl.yaml
    fi
    sudo -E ip netns exec "${UPFNS}${i}" ./bin/upf -c ./config/upfcfg.testulcl.yaml &
    sleep 1
    sed -i -e "s/internet2/internet/g" ./config/upfcfg.testulcl.yaml
done

cd test
$GOROOT/bin/go test -v -vet=off -run $1

sleep 3
sudo killall -15 upf
sleep 1

cd ../..
mkdir -p testkeylog
for KEYLOG in $(ls *sslkey.log); do
    mv $KEYLOG testkeylog
done

sudo ip addr del 10.60.0.1/32 dev lo
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
