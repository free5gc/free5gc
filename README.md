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
  - [A. Prerequisites](#a-prerequisites)
  - [B. Install Control Plane Elements](#b-install-control-plane-elements)
  - [C. Install User Plane Function (UPF)](#c-install-user-plane-function-upf)
- [Run](#run)
  - [A. Run the Core Network](#a-run-the-core-network)
  - [B. Run the N3IWF (Individually)](#b-run-the-n3iwf-individually)
  - [C. Run all-in-one with outside RAN](#c-run-all-in-one-with-outside-ran)
  - [D. Deploy within containers](#d-deploy-within-containers)
- [Test](#test)
- [More information](#more-information)
- [Release Note](#release-note)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Hardware Tested

In the market today, there are neither hardware gNB nor hardware UEs that interface directly with a standalone 5GC, so no hardware testing has yet been performed against free5gc. 

## Questions

For questions and support please use the [official forum](https://forum.free5gc.org). The issue list of this repo is exclusively for bug reports and feature requests.

## Recommended Environment

free5gc has been tested against the following environment:

- Software
    - OS: Ubuntu 18.04
    - gcc 7.3.0
    - Go 1.14.4 linux/amd64
    - kernel version 5.0.0-23-generic

The listed kernel version is required for the UPF element.

- Minimum Hardware
    - CPU: Intel i5 processor
    - RAM: 4GB
    - Hard drive: 160GB
    - NIC: Any 1Gbps Ethernet card supported in the Linux kernel

- Recommended Hardware
    - CPU: Intel i7 processor
    - RAM: 8GB
    - Hard drive: 160GB
    - NIC: Any 10Gbps Ethernet card supported in the Linux kernel

This guide assumes that you will run all 5GC elements on a single machine.

## Installation

### A. Prerequisites

1. Linux Kernel Version
    * In order to use the UPF element, you must use the `5.0.0-23-generic` version of the Linux kernel.  free5gc uses the [gtp5g kernel module](https://github.com/PrinzOwO/gtp5g), which has been tested and compiled against that kernel version only.  To determine the version of the Linux kernel you are using:

    ```bash
        $ uname -r
        5.0.0-23-generic
    ```

You will not be able to run most of the tests in [Test](#test) section unless you deploy a UPF.

2. Golang Version
    * As noted above, free5gc is built and tested with Go 1.14.4
    * To check the version of Go on your system, from a command prompt:

    ```bash
        go version
    ```

    * If another version of Go is installed, remove the existing version and install Go 1.14.4:

    ```bash
        # this assumes your current version of Go is in the default location
        sudo rm -rf /usr/local/go
        wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
        sudo tar -C /usr/local -zxvf go1.14.4.linux-amd64.tar.gz
    ```

    * If Go is not installed on your system:

    ```bash
        wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
        sudo tar -C /usr/local -zxvf go1.14.4.linux-amd64.tar.gz
        mkdir -p ~/go/{bin,pkg,src}
        # The following assume that your shell is bash
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
        source ~/.bashrc
    ```

    * Further information and installation instructions for `golang` are available at the [official golang site](https://golang.org/doc/install).

3. Control-plane Supporting Pacakges

```bash
sudo apt -y update
sudo apt -y install mongodb wget git
sudo systemctl start mongodb
```

4. User-plane Supporting Packages

```bash
sudo apt -y update
sudo apt -y install git gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
go get -u github.com/sirupsen/logrus
```

5. Linux Host Network Settings

```bash
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -t nat -A POSTROUTING -o <dn_interface> -j MASQUERADE
sudo systemctl stop ufw
```

### B. Install Control Plane Elements
    
1. Clone the free5GC repository
    * To install the latest stable build (v3.0.4):

    ```bash
        cd ~
        git clone --recursive -b v3.0.4 -j `nproc` https://github.com/free5gc/free5gc.git
        cd free5gc
    ```

    * Alternatively, if you wish to install the latest nightly build:

    ```bash
        cd ~/free5gc
        git checkout master
        git submodule sync
        git submodule update --init --jobs `nproc`
        git submodule foreach git checkout master
        git submodule foreach git pull --jobs `nproc`
    ```

2. Install all Go module dependencies

```bash
cd ~/free5gc
go mod download
```
**NOTE: the root folder name for this repository must be `free5gc`.  If it is changed, compilation will fail.**

3. Compile network function services in `free5gc`
    * To do so individually (e.g., AMF only):

    ```bash
        cd ~/free5gc
        go build -o bin/amf -x src/amf/amf.go
    ```

    * To build all network functions:

    ```bash
        cd ~/free5gc
        ./build.sh
    ```

### C. Install User Plane Function (UPF)
    
1. As noted above, the GTP kernel module used by the UPF requires that you use Linux kernel version `5.0.0-23-generic`.  To verify your version:

```bash
uname -r
```

2. Retrieve the 5G GTP-U kernel module using `git` and build it

```bash
git clone -b v0.2.0 https://github.com/PrinzOwO/gtp5g.git
cd gtp5g
make
sudo make install
```

3. Build the UPF (you may skip this step if you built all network functions above):

   a. to build using make:
   
   ```bash
   cd ~/free5gc
   make upf
   ```
  
   b. alternatively, to build manually:

   ```bash
   cd ~/free5gc/src/upf
   mkdir build
   cd build
   cmake ..
   make -j`nproc`
   ```

4. Customize the UPF as desired.  The UPF configuration file is `free5gc/src/upf/build/config/upfcfg.yaml`.

## Run

### A. Run the Core Network 

Option 1. Run network function services individually.  For example, to run the AMF:

```bash
cd ~/free5gc
./bin/amf
```

**Note: The N3IWF needs specific configuration, which is detailed in section B.** 

Option 2. Run whole core network

```bash
cd ~/free5gc
./run.sh
```

### B. Run the N3IWF (Individually)

To run an instance of the N3IWF, make sure your system is equipped with three network interfaces: the first connects to the AMF, the second connects to the UPF, and the third is for IKE daemon.

Configure each interface with a suitable IP address.

Create an interface for IPSec traffic:

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

Run the N3IWF (root privilege is required):

```bash
cd ~/free5gc/
sudo ./bin/n3iwf
```

### C. Run all-in-one with external RAN

Refer to this [sample config](./sample/ran_attach_config) if you wish to connect an external RAN with a complete free5GC core network.

### D. Deploy within containers

[free5gc-compose](https://github.com/free5gc/free5gc-compose/) provides a sample for the deployment of elements within containers.

## Test

Start a Wireshark capture on any core-connected interface, applying the filter `'pfcp||icmp||gtp'`.

In order to run the tests, first do this:

```bash
cd ~/free5gc
chmod +x ./test.sh
 ```

The tests are all run from within `~/free5gc`.

a. TestRegistration

```bash
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

```bash
./test.sh TestPaging
```

h. TestN2Handover

```bash
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

## More information

For more details, reference the free5gc [wiki](https://github.com/free5gc/free5gc/wiki).

## Release Note

Detailed changes for each release are documented in the [release notes](https://github.com/free5gc/free5gc/releases).
