# free5GC Stage 2 Installation Guide
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffree5gc%2Ffree5gc-stage-2.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffree5gc%2Ffree5gc-stage-2?ref=badge_shield)



## Minimum Requirement
- Software
    - OS: Ubuntu 18.04 or later versions
    - gcc 7.3.0
    - Go 1.12.9 linux/amd64
    - QEMU emulator 2.11.1
```bash
**Note:** Please use Ubuntu 18.04 or later versions and go 1.12.9 linux/amd64
```

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

## Hardware Tested 
There are no gNB and UE for standalone 5GC available in the market yet.


## Installation

### A. Install Control Plane Entities

1. Install the required packages
    - ```sudo apt -y update```
    - ```sudo apt -y install mongodb wget git```
    - ```sudo systemctl start mongodb```
2. Go installation
    * If another version of Go is installed
        - Please remove the previous Go version
            - ```sudo rm -rf /usr/local/go```
        - Install Go 1.12.9
            - ```wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz```
            - ```sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz```
    * Clean installation
        - Install Go 1.12.9
            - ```wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz```
            - ```sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz```
            - ```mkdir -p ~/go/{bin,pkg,src}```
            - ```echo 'export GOPATH=$HOME/go' >> ~/.bashrc```
            - ```echo 'export GOROOT=/usr/local/go' >> ~/.bashrc```
            - ```echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc```
            - ```echo 'export GO111MODULE=off' >> ~/.bashrc```
            - ```source ~/.bashrc```

3. Clone free5GC project in `$GOPATH/src`
    - ```cd $GOPATH/src```
    - ```git clone https://bitbucket.org/free5GC/free5gc-stage-2.git free5gc```
4. Run the script to install dependent packages
    - ```cd $GOPATH/src/free5gc```
    - ```chmod +x ./install_env.sh```
    - ```./install_env.sh```
    
    - ```Please ignore error messages during the package dependencies installation process.```

5. Extract the `free5gc_libs.tar.gz` to setup the environment for compiling
    - ```cd $GOPATH/src/free5gc```
    - ```tar -C $GOPATH -zxvf free5gc_libs.tar.gz```
6. Compile network function services in `$GOPATH/src/free5gc`, e.g. AMF:
    - ```cd $GOPATH/src/free5gc```
    - ```go build -o bin/amf -x src/amf/amf.go```
7. Run network function services, e.g. AMF:
    - ```cd $GOPATH/src/free5gc```
    - ```./bin/amf```

    - ```In step 3, the folder name should remain free5gc. Please do not modify it or the compilation would fail.```


### B. Install User Plane Entity (UPF)
1. Install the required packages
    ```bash
    sudo apt -y update
    sudo apt -y install git gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
    go get -u github.com/sirupsen/logrus
    ```
2. Enter the UPF directory
    - ```cd $GOPATH/src/free5gc/src/upf```
3. Build from sources
    - ```mkdir build```
    - ```cd build```
    - ```cmake ..```
    - ```make -j `nproc` ```
4. Run UPF library test
    - ```(In directory: $GOPATH/src/free5gc/src/upf/build)```
    - ```sudo ./bin/testgtpv1```
5. Config is located at `$GOPATH/src/free5gc/src/upf/build/config/upfcfg.yaml`

### C. Run Procedure Tests
Start Wireshark to capture any interface with pfcp||icmp||gtp filter and run the tests below to simulate the procedures:
```bash
cd $GOPATH/src/free5gc
chmod +x ./test.sh
```
a. TestRegistration
```bash
(In directory: $GOPATH/src/free5gc)
sudo ./test.sh TestRegistration
```
b. TestServiceRequest
```bash
sudo ./test.sh TestServiceRequest
```
c. TestXnHandover
```bash
sudo ./test.sh TestXnHandover
```
d. TestDeregistration
```bash
sudo ./test.sh TestDeregistration
```
e. TestPDUSessionReleaseRequest
```bash
sudo ./test.sh TestPDUSessionReleaseRequest
```

### Appendix A: System Environment Cleaning
The below commands may be helpful for development purposes.

1. Remove POSIX message queues
    - ```ls /dev/mqueue/```
    - ```rm /dev/mqueue/*```
2. Remove gtp tunnels (using tools in libgtpnl)
    - ```cd ./src/upf/lib/libgtpnl-1.2.1/tools```
    - ```./gtp-tunnel list```
3. Remove gtp devices (using tools in libgtpnl)
    - ```cd ./src/upf/lib/libgtpnl-1.2.1/tools```
    - ```sudo ./gtp-link del {Dev-Name}```
## Appendix B: Program the SIM Card
Install packages:
```bash
sudo apt-get install pcscd pcsc-tools libccid python-dev swig python-setuptools python-pip libpcsclite-dev
sudo pip install pycrypto
```

Download PySIM
```bash
git clone git://git.osmocom.org/pysim.git
```

Change to pyscard folder and install
```bash
cd <pyscard-path>
sudo /usr/bin/python setup.py build_ext install
```

Verify your reader is ready

```bash
sudo pcsc_scan
```

Check whether your reader can read the SIM card
```bash
cd <pysim-path>
./pySim-read.py –p 0
```

Program your SIM card information
```bash
./pySim-prog.py -p 0 -x 208 -y 93 -t sysmoUSIM-SJS1 -i 208930000000003 --op=8e27b6af0e692e750f32667a3b14605d -k 8baf473f2f8fd09487cccbd7097c6862 -s 8988211000000088313 -a 23605945
```

You can get your SIM card from [**sysmocom**](http://shop.sysmocom.de/products/sysmousim-sjs1-4ff). You also need a card reader to write your SIM card. You can get a card reader from [**here**](https://24h.pchome.com.tw/prod/DCAD59-A9009N6WF) or use other similar devices.

# Release Note
## v2.0.1
+ Add buffering and paging

## v2.0.2
+ Add handover feature
+ Add webui


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffree5gc%2Ffree5gc-stage-2.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffree5gc%2Ffree5gc-stage-2?ref=badge_large)