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

while getopts 'o' OPT;
do
    case $OPT in
        o) DUMP_NS=True;;
    esac
done
shift $(($OPTIND - 1))

TEST_POOL="TestRegistration|TestGUTIRegistration|TestServiceRequest|TestXnHandover|TestN2Handover|TestDeregistration|TestPDUSessionReleaseRequest|TestPaging|TestNon3GPP|TestReSynchronisation"
if [[ ! "$1" =~ $TEST_POOL ]]
then
    echo "Usage: $0 [ ${TEST_POOL//|/ | } ]"
    exit 1
fi

cp config/test/smfcfg.single.test.conf config/test/smfcfg.test.conf

GOPATH=$HOME/go
if [ $OS == "Ubuntu" ]; then
    GOROOT=/usr/local/go
elif [ $OS == "Fedora" ]; then
    GOROOT=/usr/lib/golang
fi
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

UPFNS="UPFns"
EXEC_UPFNS="sudo -E ip netns exec ${UPFNS}"

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
${EXEC_UPFNS} ip addr add 60.60.0.101 dev lo
${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1

if [ ${DUMP_NS} ]
then
    ${EXEC_UPFNS} tcpdump -i any -w ${UPFNS}.pcap &
    TCPDUMP_PID=$(sudo ip netns pids ${UPFNS})
    sudo -E tcpdump -i lo -w default_ns.pcap &
    LOCALDUMP=$!
fi

cd src/upf/build && ${EXEC_UPFNS} ./bin/free5gc-upfd -f config/upfcfg.test.yaml &
sleep 2

if [[ "$1" == "TestNon3GPP" ]]
then
    UENS="UEns"
    EXEC_UENS="sudo ip netns exec ${UENS}"

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

    # Configuration
    cp -f config/amfcfg.conf config/amfcfg.conf.bak
    cp -f config/amfcfg.n3test.conf config/amfcfg.conf

    # Run CN
    cd src/test && $GOROOT/bin/go test -v -vet=off -timeout 0 -run TestCN &
    sleep 10

    # Run N3IWF
    cd src/n3iwf && sudo -E $GOROOT/bin/go run n3iwf.go &
    sleep 5

    # Run Test UE
    cd src/test
    ${EXEC_UENS} $GOROOT/bin/go test -v -vet=off -timeout 0 -run TestNon3GPPUE -args noinit

else
    cd src/test
    $GOROOT/bin/go test -v -vet=off -run $1
fi

sleep 3
sudo killall -15 free5gc-upfd
sleep 1

if [ ${DUMP_NS} ]
then
    ${EXEC_UPFNS} kill -SIGINT ${TCPDUMP_PID}
    sudo -E kill -SIGINT ${LOCALDUMP}
fi

cd ../..
mkdir -p testkeylog
for KEYLOG in $(ls *sslkey.log); do
    mv $KEYLOG testkeylog
done

sudo ip link del veth0
sudo ip netns del ${UPFNS}
sudo ip addr del 60.60.0.1/32 dev lo

if [[ "$1" == "TestNon3GPP" ]]
then
    sudo ip xfrm policy flush
    sudo ip xfrm state flush
    sudo ip link del veth2
    sudo ip link del ipsec0
    ${EXEC_UENS} ip link del ipsec0
    sudo ip netns del ${UENS}
    sudo killall n3iwf
    killall test.test
    cp -f config/amfcfg.conf.bak config/amfcfg.conf
    rm -f config/amfcfg.conf.bak
fi

sleep 2
