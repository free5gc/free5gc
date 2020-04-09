# free5GC Stage 3 Installation Guide


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
    ```bash
    sudo apt -y update
    sudo apt -y install mongodb wget git
    sudo systemctl start mongodb
    ```
2. Go installation
    * If another version of Go is installed
        - Please remove the previous Go version
            - ```sudo rm -rf /usr/local/go```
        - Install Go 1.12.9
            ```bash
            wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz
            ```
    * Clean installation
        - Install Go 1.12.9
             ```bash
            wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz
            mkdir -p ~/go/{bin,pkg,src}
            echo 'export GOPATH=$HOME/go' >> ~/.bashrc
            echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
            echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
            echo 'export GO111MODULE=off' >> ~/.bashrc
            source ~/.bashrc
            ```

3. Clone free5GC project in `$GOPATH/src`
    ```bash
    cd $GOPATH/src
    git clone https://bitbucket.org/free5GC/free5gc-stage-2.git free5gc
    ```
4. Run the script to install dependent packages
    ```bash
    cd $GOPATH/src/free5gc
    chmod +x ./install_env.sh```
    ./install_env.sh
    
    Please ignore error messages during the package dependencies installation process.
    ```

5. Extract the `free5gc_libs.tar.gz` to setup the environment for compiling
    ```bash
    cd $GOPATH/src/free5gc
    tar -C $GOPATH -zxvf free5gc_libs.tar.gz
    ```
6. Compile network function services in `$GOPATH/src/free5gc`, e.g. AMF:
    ```bash
    cd $GOPATH/src/free5gc
    go build -o bin/amf -x src/amf/amf.go
7. Run network function services, e.g. AMF:
    ```bash
    cd $GOPATH/src/free5gc
    ./bin/amf

    In step 3, the folder name should remain free5gc. Please do not modify it or the compilation would fail.
    ```


### B. Install User Plane Entity (UPF)
1. Install the required packages
    ```bash
    sudo apt -y update
    sudo apt -y install git gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
    go get -u github.com/sirupsen/logrus
    ```
2. Please check Linux kernel version if it is 5.0.0-23-generic
    ```bash
    uname -r
    ```


    Get Linux kernel module 5G GTP-U
    ```bash
    git clone https://github.com/PrinzOwO/gtp5g.git
    cd gtp5g
    make
    sudo make install
    ```
3. Enter the UPF directory
    ```cd $GOPATH/src/free5gc/src/upf```
4. Build from sources
    ```bash
    mkdir build
    cd build
    cmake ..
    make -j `nproc`
    ```
5. Run UPF library test
    ```bash
    (In directory: $GOPATH/src/free5gc/src/upf/build
    
    sudo ./bin/testgtpv1
    ```
6. Config is located at `$GOPATH/src/free5gc/src/upf/build/config/upfcfg.yaml`



## Run

### A. Run Core Network
```
./run.sh
```

### B. Run N3IWF (manually)
To run N3IWF, make sure the machine is equipped with three network interfaces, one is for connecting AMF, another is for connecting UPF, the other is for IKE daemon.

For each interface, configure it with a suitable IP address.

We have to create an interface for IPSec traffic:
```bash
# replace <...> to suitable value
sudo ip link add ipsec0 type vti local <IKEBindAddress> remote 0.0.0.0 key <IPSecInterfaceMark>
```
Give an address to this interface, then turn on it:
```bash
# replace <...> to suitable value
sudo ip address add <IPSecInterfaceAddress/CIDRPrefix> dev ipsec0
sudo ip link set dev ipsec0 up
```

Run N3IWF (need root privilege):
```bash
cd bin/
sudo ./n3iwf
```

### C. Run Procedure Tests
Start Wireshark to capture any interface with `pfcp||icmp||gtp` filter and run the tests below to simulate the procedures:
```bash
cd $GOPATH/src/free5gc
chmod +x ./test.sh
```
a. TestRegistration
```bash
(In directory: $GOPATH/src/free5gc)
./test.sh TestRegistration
```
b. TestServiceRequest
```bash
./test.sh TestServiceRequest
```
c. TestXnHandover
```bash
./test.sh TestXnHandover
```
d. TestDeregistration
```bash
./test.sh TestDeregistration
```
e. TestPDUSessionReleaseRequest
```bash
./test.sh TestPDUSessionReleaseRequest
```

g. TestNon3GPP
```bash
./test.sh TestNon3GPP
```

h. TestULCL
```bash
./test_ulcl.sh -om 3 TestRegistration
```

## Appendix A: OAM 
1. Run the web console server
```
cd webconsole
go run server.go
```
2. Access the web console by
```
URL: http://localhost:5000
Username: admin
Password: free5gc
```
3. Now you can see the information of currently registered UEs (e.g. Supi, connected state, etc.) in the core network at the tab "DASHBOARD" of free5GC webconsole

**Note: You can add the subscribers here too**

## Appendix B: System Environment Cleaning
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

## Appendix C: Change Kernel Version
1. Check the previous kernel version: `uname -r`
2. search specific kernel version and install, take `5.0.0-23-generic` for example
```bash
sudo apt search 'linux-image-5.0.0-23-generic'
sudo apt install 'linux-image-5.0.0-23-generic'
sudo apt install 'linux-headers-5.0.0-23-generic'
```
3. Update initramfs and grub
```bash
sudo update-initramfs -u -k all
sudo update-grub
```
4. Reboot, enter grub and choose kernel version `5.0.0-23-generic`
```bash
sudo reboot
```
#### Optional: Remove Kernel Image
```
sudo apt remove 'linux-image-5.0.0-23-generic'
sudo apt remove 'linux-headers-5.0.0-23-generic'
```

## Appendix D: Program the SIM Card
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


## Trouble Shooting

1. `ERROR: [SCTP] Failed to connect given AMF    N3IWF=NGAP`

    This error occured when N3IWF was started before AMF finishing initialization. This error usually appears when you run the TestNon3GPP in the first time.

    Rerun the test should be fine. If it still not be solved, larger the sleeping time in line 110 of `test.sh`.

2. TestNon3GPP will modify the `config/amfcfg.conf`. So, if you had killed the TestNon3GPP test before it finished, you might need to copy `config/amfcfg.conf.bak` back to `config/amfcfg.conf` to let other tests pass.

    `cp config/amfcfg.conf.bak config/amfcfg.conf`

# Release Note
## v3.0.0
+ AMF
    + Support SMF selection at PDU session establishment
    + Fix SUCI handling procedure
+ SMF
    + Feature
        + ULCL by config
        + Authorized QoS
    + Bugfix
        + PDU Session Establishment PDUAddress Information
        + PDU Session Establishment N1 Message
        + SMContext Release Procedure
+ UPF:
    + ULCL feature
    + support SDF Filter
    + support N9 interface
+ OAM
    + Get Registered UE Context API
    + OAM web UI to display Registered UE Context
+ N3IWF
    + Support Registration procedure for untrusted non-3GPP access
    + Support UE Requested PDU Session Establishment via Untrusted non-3GPP Access
+ UDM
    + SUCI to SUPI de-concealment
    + Notification 
        + Callback notification to NF ( in SDM service)
        + UDM initiated deregistration notification to NF ( in UECM service)

## v2.0.2
+ Add debug mode on NFs
+ Auto add Linux routing when UPF runs
+ Add AMF consumer for AM policy
+ Add SM policy
+ Allow security NIA0 and NEA0
+ Add handover feature
+ Add webui
+ Update license
+ Bugfix for incorrect DNN
+ Bugfix for NFs registering to NRF

## v2.0.1
+ Global
    + Update license and readme
    + Add Paging feature
    + Bugfix for AN release issue
    + Add URL for SBI in NFs' config

+ AMF
    + Add Paging feature
    + Bugfix for SCTP PPID to 60
    + Bugfix for UE release in testing
    + Bugfix for too fast send UP data in testing
    + Bugfix for sync with defaultc config in testing

+ SMF
    + Add Paging feature
    + Create PDR with FAR ID
    + Bugfix for selecting DNN fail handler

+ UPF
    + Sync config default address with Go NFs
    + Remove GTP tunnel by removing PDR/FAR
    + Bugfix for PFCP association setup
    + Bugfix for new PDR/FAR creating
    + Bugfix for PFCP session report
    + Bugfix for getting from PDR
    + Bugfix for log format and update logger version

+ PCF
    + Bugfix for lost field and method

