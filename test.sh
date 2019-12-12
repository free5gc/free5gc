#!/bin/bash

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

while getopts 'o' OPT;
do
    case $OPT in
        o) DUMP_NS=True;;
    esac
done
shift $(($OPTIND - 1))

TEST_POOL="TestRegistration|TestServiceRequest|TestXnHandover|TestN2Handover|TestDeregistration|TestPDUSessionReleaseRequest|TestPaging"
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
GO111MODULE=off

UPFNS="UPFns"
EXEC_UPFNS="sudo ip netns exec ${UPFNS}"

export GIN_MODE=release

# Setup network namespace
sudo ip netns add ${UPFNS}

sudo ip link add veth0 type veth peer name veth1
sudo ip link set veth0 up
sudo ip addr add 60.60.0.1 dev lo
sudo ip addr add 10.200.200.1/24 dev veth0
sudo ip addr add 10.200.200.2/24 dev veth0

sudo ip link set veth1 netns ${UPFNS}

${EXEC_UPFNS} ip link set lo up
${EXEC_UPFNS} ip link set veth1 up
${EXEC_UPFNS} ip addr add 60.60.0.100 dev lo
${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1

if [ ${DUMP_NS} ]
then
    ${EXEC_UPFNS} tcpdump -i any -w ${UPFNS}.pcap &
    TCPDUMP_PID=$(sudo ip netns pids ${UPFNS})
fi

cd src/upf/build && ${EXEC_UPFNS} ./bin/free5gc-upfd -f config/upfcfg.test.yaml &
sleep 2

cd src/test
$GOROOT/bin/go test -v -vet=off -run $1

sleep 1
sudo killall -15 free5gc-upfd
sleep 1

if [ ${DUMP_NS} ]
then
    ${EXEC_UPFNS} kill -SIGINT ${TCPDUMP_PID}
fi

cd ../..
mkdir -p testkeylog
for KEYLOG in $(ls *sslkey.log); do 
    mv $KEYLOG testkeylog
done

sudo ip link del veth0
sudo ip netns del ${UPFNS}
sudo ip addr del 60.60.0.1/32 dev lo
