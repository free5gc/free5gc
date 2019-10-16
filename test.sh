#!/bin/bash

if [ -z "$1" ]
then
    echo "Usage: $0 [TestRegistration | TestServiceRequest | TestXnHandover | TestDeregistration | TestPDUSessionReleaseRequest]"
    exit 1
fi

GOPATH=$HOME/go
GOROOT=/usr/local/go
PATH=$PATH:$GOPATH/bin:$GOROOT/bin
GO111MODULE=off

UPFNS="UPFns"
EXEC_UPFNS="ip netns exec ${UPFNS}"

export GIN_MODE=release

# Setup network namespace
ip netns add ${UPFNS}

ip link add veth0 type veth peer name veth1
ip link set veth0 up
ip addr add 60.60.0.1 dev lo
ip addr add 10.200.200.1/24 dev veth0

ip link set veth1 netns ${UPFNS}

${EXEC_UPFNS} ip link set lo up
${EXEC_UPFNS} ip link set veth1 up
${EXEC_UPFNS} ip addr add 60.60.0.100 dev lo
${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1

cd src/upf/build && ${EXEC_UPFNS} ./bin/free5gc-upfd &
sleep 2
${EXEC_UPFNS} ip r add 60.60.0.0/24 dev free5GCgtp0

cd src/test
$GOROOT/bin/go test -v -vet=off -run $1

killall -15 free5gc-upfd
sleep 1

ip link del veth0
ip netns del ${UPFNS}
ip addr del 60.60.0.1/32 dev lo
