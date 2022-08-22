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

TEST_POOL="TestRegistration|TestGUTIRegistration|TestServiceRequest|TestXnHandover|TestN2Handover|TestDeregistration|TestPDUSessionReleaseRequest|TestPaging|TestNon3GPP|TestReSynchronization|TestDuplicateRegistration|TestEAPAKAPrimeAuthentication"
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
EXEC_UPFNS="sudo -E ip netns exec ${UPFNS}"

export GIN_MODE=release

function terminate()
{
    sleep 3
    sudo killall -15 upf

    if [ ${DUMP_NS} ]
    then
        # kill all tcpdump processes in the default network namespace
        sudo killall tcpdump
        sleep 1
    fi

    sudo ip link del veth0
    sudo ip netns del ${UPFNS}
    sudo ip addr del 10.60.0.1/32 dev lo

    if [[ "$1" == "TestNon3GPP" ]]
    then
        if [ ${DUMP_NS} ]
        then
            cd .. && sudo ip xfrm state > ${PCAP_PATH}/NWu_SA_state.log
        fi
        sudo ip xfrm policy flush
        sudo ip xfrm state flush
        sudo ip netns del ${UENS}
        removeN3iwfInterfaces
        sudo ip link del veth2
        sudo killall n3iwf
        killall test.test
    fi

    sleep 5
}

function removeN3iwfInterfaces()
{
    # Remove all GRE interfaces
    GREs=$(ip link show type gre | awk 'NR%2==1 {print $2}' | cut -d @ -f 1)
    for GRE in ${GREs}; do
        sudo ip link del ${GRE}
    done

    # Remove all XFRM interfaces
    XFRMIs=$(ip link show type xfrm | awk 'NR%2==1 {print $2}' | cut -d @ -f 1)
    for XFRMI in ${XFRMIs}; do
        sudo ip link del ${XFRMI}
    done
}

function handleSIGINT()
{
    echo -e "\033[41;37m Terminating due to SIGINT ... \033[0m"
    terminate $1
}

trap handleSIGINT SIGINT

function setupN3ueEnv()
{
    UENS="UEns"
    EXEC_UENS="sudo -E ip netns exec ${UENS}"

    sudo ip netns add ${UENS}

    sudo ip link add veth2 type veth peer name veth3
    sudo ip addr add 192.168.127.1/24 dev veth2
    sudo ip link set veth2 up

    sudo ip link set veth3 netns ${UENS}
    ${EXEC_UENS} ip addr add 192.168.127.2/24 dev veth3
    ${EXEC_UENS} ip link set lo up
    ${EXEC_UENS} ip link set veth3 up
    ${EXEC_UENS} ip a
}

function tcpdumpN3IWF()
{
    N3IWF_IPSec_iface_addr=192.168.127.1
    N3IWF_IPsec_inner_addr=10.0.0.1
    N3IWF_GTP_addr=10.200.200.2
    UE_DN_addr=10.60.0.1

    ${EXEC_UENS} tcpdump -U -i any -w $PCAP_PATH/$UENS.pcap &
    TCPDUMP_QUERY=" host $N3IWF_IPSec_iface_addr or \
                    host $N3IWF_IPsec_inner_addr or \
                    host $N3IWF_GTP_addr or \
                    host $UE_DN_addr"
    sudo -E tcpdump -U -i any $TCPDUMP_QUERY -w $PCAP_PATH/n3iwf.pcap &
}

# Setup network namespace
sudo ip netns add ${UPFNS}

sudo ip link add veth0 type veth peer name veth1
sudo ip link set veth0 up
sudo ip addr add 10.60.0.1 dev lo
sudo ip addr add 10.200.200.1/24 dev veth0
sudo ip addr add 10.200.200.2/24 dev veth0

sudo ip link set veth1 netns ${UPFNS}

${EXEC_UPFNS} ip link set lo up
${EXEC_UPFNS} ip link set veth1 up
${EXEC_UPFNS} ip addr add 10.60.0.101 dev lo
${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1

if [ ${DUMP_NS} ]
then
    PCAP_PATH=testpcap
    mkdir -p ${PCAP_PATH}
    ${EXEC_UPFNS} tcpdump -U -i any -w ${PCAP_PATH}/${UPFNS}.pcap &
    sudo -E tcpdump -U -i lo -w ${PCAP_PATH}/default_ns.pcap &
fi

${EXEC_UPFNS} ./bin/upf -c ./config/upfcfg.test.yaml &
sleep 2

if [[ "$1" == "TestNon3GPP" ]]
then
    removeN3iwfInterfaces
    # setup N3UE's namespace, interfaces for IPsec
    setupN3ueEnv
    if [ ${DUMP_NS} ]
    then
        tcpdumpN3IWF
    fi

    # Run CN
    cd test && go test -v -vet=off -timeout 0 -run TestCN &
    sleep 10

    # Run N3IWF
    sudo -E ./bin/n3iwf -c ./config/n3iwfcfg.test.yaml &
    sleep 5

    # Run Test UE
    cd test
    if ! go test -v -vet=off -timeout 0 -run TestNon3GPPUE -args noinit; then
        echo "Test result: Failed"
        terminate $1
        exit 1
    else
        echo "Test result: Succeeded"
    fi
else
    cd test
    if ! go test -v -vet=off -run $1; then
        echo "Test result: Failed"
        terminate $1
        exit 1
    else
        echo "Test result: Succeeded"
    fi
fi

terminate $1
