<p align="center">
<a href="https://free5gc.org"><img width="40%" src="https://forum.free5gc.org/uploads/default/original/1X/324695bfc6481bd556c11018f2834086cf5ec645.png" alt="free5GC"/></a>
</p>

<p align="center">
<a href="https://github.com/free5gc/free5gc/releases"><img src="https://img.shields.io/github/v/release/free5gc/free5gc?color=orange" alt="Release"/></a>
<a href="https://github.com/free5gc/free5gc/blob/master/LICENSE.txt"><img src="https://img.shields.io/github/license/free5gc/free5gc?color=blue" alt="License"/></a>
<a href="https://forum.free5gc.org"><img src="https://img.shields.io/discourse/topics?server=https%3A%2F%2Fforum.free5gc.org&color=lightblue" alt="Forum"/></a>
<a href="https://www.codefactor.io/repository/github/free5gc/free5gc"><img src="https://www.codefactor.io/repository/github/free5gc/free5gc/badge" alt="CodeFactor" /></a>
<a href="https://goreportcard.com/report/github.com/free5gc/free5gc"><img src="https://goreportcard.com/badge/github.com/free5gc/free5gc" alt="Go Report Card" /></a>
<a href="https://github.com/free5gc/free5gc/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen" alt="PRs Welcome"/></a>
</p>


## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Hardware Tested](#hardware-tested)
- [Questions](#questions)
- [Recommended Environment](#recommended-environment)
- [Installation](#installation)
  - [A. Pre-requisite](#a-pre-requisite)
  - [B. Install Control Plane Entities](#b-install-control-plane-entities)
  - [C. Install User Plane Function (UPF)](#c-install-user-plane-function-upf)
- [Run](#run)
  - [A. Run Core Network](#a-run-core-network)
  - [B. Run N3IWF (Individually)](#b-run-n3iwf-individually)
  - [C. Run all in one with outside RAN](#c-run-all-in-one-with-outside-ran)
  - [D. Deploy with container](#d-deploy-with-container)
- [Test](#test)
- [Release Note](#release-note)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Hardware Tested
There are no gNB and UE for standalone 5GC available in the market yet.

## Questions
For questions and support please use the [official forum](https://forum.free5gc.org). The issue list of this repo is exclusively
for bug reports and feature requests.

## Recommended Environment
- Software
    - OS: Ubuntu 18.04
    - gcc 7.3.0
    - Go 1.14.4 linux/amd64
    - kernel version 5.0.0-23-generic (MUST for UPF)

**Note: Please use Ubuntu 18.04 and kernel version 5.0.0-23-generic**


You can use `go version` to check your current Go version.
```bash
- Hardware
    - CPU: Intel i5 processor
    - RAM: 4GB
    - Hard drive: 160G
    - NIC card: 1Gbps ethernet card

- Hardware recommended
    - CPU: Intel i7 processor
    - RAM: 8GB
    - Hard drive: 160G
    - NIC card: 10Gbps ethernet card
```


## Installation
### A. Pre-requisite

0. Required kernel version `5.0.0-23-generic`. This request is from the module
   [gtp5g](https://github.com/PrinzOwO/gtp5g) that we has used. Any more details
   please check [here](https://github.com/PrinzOwO/gtp5g)
   ```bash
   # Check kernel version
   $ uname -r
   5.0.0-23-generic
   ```

1. Require go language
    * If another version of Go is installed
        - Please remove the previous Go version
            - `sudo rm -rf /usr/local/go`
        - Install Go 1.14.4
            ```bash
            wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.14.4.linux-amd64.tar.gz
            ```
    * Clean installation
        - Install Go 1.14.4
             ```bash
            wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.14.4.linux-amd64.tar.gz
            mkdir -p ~/go/{bin,pkg,src}
            echo 'export GOPATH=$HOME/go' >> ~/.bashrc
            echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
            echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
            source ~/.bashrc
            ```

2. Required packages for control plane
    ```bash
    sudo apt -y update
    sudo apt -y install mongodb wget git
    sudo systemctl start mongodb
    ```

3. Required packages for user plane
    ```bash
    sudo apt -y update
    sudo apt -y install git gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
    go get -u github.com/sirupsen/logrus
    ```

4. Network Setting
    ```bash
    sudo sysctl -w net.ipv4.ip_forward=1
    sudo iptables -t nat -A POSTROUTING -o <dn_interface> -j MASQUERADE
    sudo systemctl stop ufw
    ```

### B. Install Control Plane Entities

1. Clone free5GC project
    ```bash
    cd ~
    git clone --recursive -b v3.0.4 -j `nproc` https://github.com/free5gc/free5gc.git
    cd free5gc
    ```

    (Optional) If you want to use the nightly version, runs:
    ```bash
    cd ~/free5gc
    git checkout master
    git submodule sync
    git submodule update --init --jobs `nproc`
    git submodule foreach git checkout master
    git submodule foreach git pull --jobs `nproc`
    ```

2. Run the script to install dependent packages
    ```bash
    cd ~/free5gc
    go mod download
    ```
    **In step 2, the folder name should remain free5gc. Please do not modify it or the compilation would fail.**

3. Compile network function services in `free5gc` individually, e.g. AMF (redo this step for each NF), or
    ```bash
    cd ~/free5gc
    make amf
    ```
    **To build all network functions in one command**
    ```bash
    cd ~/free5gc
    make all
    ```


### C. Install User Plane Function (UPF)

1. Please check Linux kernel version if it is `5.0.0-23-generic`
    ```bash
    uname -r
    ```


    Get Linux kernel module 5G GTP-U
    ```bash
    git clone -b v0.2.0 https://github.com/PrinzOwO/gtp5g.git
    cd gtp5g
    make
    sudo make install
    ```

2. Build from sources (skip this step if you run make all previously) via make, or
    ```bash
    cd ~/free5gc
    make upf
    ```
    build manually
    ```bash
    cd ~/free5gc/src/upf
    mkdir build
    cd build
    cmake ..
    make -j`nproc`
    ```

**Note: UPF's config is located at** `free5gc/src/upf/build/config/upfcfg.yaml`

## Run

### A. Run Core Network
Option 1. Run network function service individually, e.g. AMF (redo this for each NF), or
```bash
cd ~/free5gc
./bin/amf
```

**Note: For N3IWF needs specific configuration in section B**

Option 2. Run whole core network with command
```
./run.sh
```

### B. Run N3IWF (Individually)
To run N3IWF, make sure the machine is equipped with three network interfaces. (one is for connecting AMF, another is for connecting UPF, the other is for IKE daemon)

We need to configure each interface with a suitable IP address.

We have to create an interface for IPSec traffic:
```bash
# replace <...> to suitable value
sudo ip link add ipsec0 type vti local <IKEBindAddress> remote 0.0.0.0 key <IPSecInterfaceMark>
```
Assign an address to this interface, then bring it up:
```bash
# replace <...> to suitable value
sudo ip address add <IPSecInterfaceAddress/CIDRPrefix> dev ipsec0
sudo ip link set dev ipsec0 up
```

Run N3IWF (root privilege is required):
```bash
cd ~/free5gc/
sudo ./bin/n3iwf
```

### C. Run all in one with outside RAN

Reference to [sample config](./sample/ran_attach_config) if need to connect the outside RAN with all in one free5GC core network.

### D. Deploy with container

Reference to [free5gc-compose](https://github.com/free5gc/free5gc-compose/) as the sample for container deployment.

## Test
Start Wireshark to capture any interface with `pfcp||icmp||gtp` filter and run the tests below to simulate the procedures:
```bash
cd ~/free5gc
chmod +x ./test.sh
```
a. TestRegistration
```bash
(In directory: ~/free5gc)
./test.sh TestRegistration
```

b. TestGUTIRegistration
```bash
./test.sh TestGUTIRegistration
```

c. TestServiceRequest
```bash
./test.sh TestServiceRequest
```

d. TestXnHandover
```bash
./test.sh TestXnHandover
```

e. TestDeregistration
```bash
./test.sh TestDeregistration
```

f. TestPDUSessionReleaseRequest
```bash
./test.sh TestPDUSessionReleaseRequest
```

g. TestPaging
```!
./test.sh TestPaging
```

h. TestN2Handover
```!
./test.sh TestN2Handover
```

i. TestNon3GPP
```bash
./test.sh TestNon3GPP
```

j. TestReSynchronisation
```bash
./test.sh TestReSynchronisation
```

k. TestULCL
```bash
./test_ulcl.sh -om 3 TestRegistration
```

**For more details, you can reference to our [wiki](https://github.com/free5gc/free5gc/wiki)**

## Release Note
Detailed changes for each release are documented in the [release notes](https://github.com/free5gc/free5gc/releases).

